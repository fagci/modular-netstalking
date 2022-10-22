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
	port       string
	workers    int
	count      uint64
	timeout    time.Duration
	printMutex sync.Mutex
)

func init() {
	flag.StringVar(&port, "p", "80", "port")
	flag.IntVar(&workers, "w", 256, "workers count")
	flag.Uint64Var(&count, "c", 0, "count of IPs to gather (0 = infinite)")
	flag.DurationVar(&timeout, "t", 750*time.Millisecond, "connection timeout")
}

func CheckPort(host string, port string, ch chan string) {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout); err == nil {
		conn.Close()
		ch <- host
	}
}

func Work(ch_in chan string, ch_out chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if ip, ok := <-ch_in; ok {
			CheckPort(ip, port, ch_out)
			continue
		}
		return
	}
}

func Gather(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	limited := count > 0
	for {
		ip, ok := <-ch
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

func ReadStdin(ch_in chan string, ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			ch_in <- scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	close(ch_in)
}

func RunWorkers(ch_in, ch_out chan string) {
	var wg_workers sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg_workers.Add(1)
		go Work(ch_in, ch_out, &wg_workers)
	}
	wg_workers.Wait()
	close(ch_out)
}

func main() {
	var wg_gather sync.WaitGroup

	flag.Parse()

	ch_in := make(chan string, workers*2)
	ch_out := make(chan string)

	reader_ctx, reader_cancel := context.WithCancel(context.Background())

	wg_gather.Add(1)

	go Gather(ch_out, &wg_gather)
	go RunWorkers(ch_in, ch_out)
	go ReadStdin(ch_in, reader_ctx)

	wg_gather.Wait()
	reader_cancel()
}
