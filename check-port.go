package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var (
	port    string
	workers int
	count   uint64
	timeout time.Duration
)

func init() {
	flag.StringVar(&port, "p", "80", "port")
	flag.IntVar(&workers, "w", 256, "workers count")
	flag.Uint64Var(&count, "c", 0, "count of IPs to gather (0 = infinite)")
	flag.DurationVar(&timeout, "t", 750*time.Millisecond, "connection timeout")
}

func CheckPort(host string, port string, toGatherIPsChan chan<- string) {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout); err == nil {
		conn.Close()
		toGatherIPsChan <- host
	}
}

func Work(toCheckIPsChan <-chan string, toGatherIPsChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if ip, ok := <-toCheckIPsChan; ok {
			CheckPort(ip, port, toGatherIPsChan)
			continue
		}
		return
	}
}

func Gather(toGatherIPsChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	limited := count > 0
	for {
		ip, ok := <-toGatherIPsChan
		if !ok {
			break
		}
		fmt.Println(ip)
		if limited {
			count--
			if count == 0 {
				break
			}
		}
	}
}

func ReadStdin(toCheckIPsChan chan<- string, ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			toCheckIPsChan <- scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	close(toCheckIPsChan)
}

func RunWorkers(toCheckIPsChan <-chan string, toGatherIPsChan chan<- string) {
	var wg_workers sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg_workers.Add(1)
		go Work(toCheckIPsChan, toGatherIPsChan, &wg_workers)
	}
	wg_workers.Wait()
	close(toGatherIPsChan)
}

func main() {
	var wgGather sync.WaitGroup

	flag.Parse()

	toCheckIPsChan := make(chan string, workers*2)
	toGatherIPsChan := make(chan string)

	readerCtx, readerCancel := context.WithCancel(context.Background())

	wgGather.Add(1)

	go Gather(toGatherIPsChan, &wgGather)
	go RunWorkers(toCheckIPsChan, toGatherIPsChan)
	go ReadStdin(toCheckIPsChan, readerCtx)

	wgGather.Wait()
	readerCancel()
}
