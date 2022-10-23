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

	"modular-netstalking/lib"
)

var (
	port    string
	workers int
	timeout time.Duration
)

func init() {
	flag.StringVar(&port, "p", "", "port")
	flag.IntVar(&workers, "w", 256, "workers count")
	flag.DurationVar(&timeout, "t", 750*time.Millisecond, "connection timeout")
	flag.Uint64Var(&lib.OutputCount, "c", 0, "count of IPs to gather (0 = infinite)")
}

func CheckPort(host string, toGatherIPsChan chan<- lib.HostInfo) {
	if conn, err := net.DialTimeout("tcp", host, timeout); err == nil {
		conn.Close()
		toGatherIPsChan <- lib.NewHostInfo(host)
	}
}

func Work(toCheckIPsChan <-chan string, toGatherIPsChan chan<- lib.HostInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if host, ok := <-toCheckIPsChan; ok {
			CheckPort(host, toGatherIPsChan)
			continue
		}
		return
	}
}

func Gather(toGatherIPsChan <-chan lib.HostInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	limited := lib.OutputCount > 0
	rest := lib.OutputCount
	for {
		ip, ok := <-toGatherIPsChan
		if !ok {
			break
		}
		fmt.Println(ip.String(!lib.GreppableOutput))
		if limited {
			rest--
			if rest == 0 {
				break
			}
		}
	}
}

func ReadStdin(toCheckIPsChan chan<- string, ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	noPortSpecified := port == ""
	if noPortSpecified {
		scanner.Scan()
		host := scanner.Text()
		if _, _, err := net.SplitHostPort(host); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(255)
		}
		toCheckIPsChan <- host
	}
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			if noPortSpecified {
				toCheckIPsChan <- scanner.Text()
			} else {
				toCheckIPsChan <- net.JoinHostPort(scanner.Text(), port)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	close(toCheckIPsChan)
}

func RunWorkers(toCheckIPsChan <-chan string, toGatherIPsChan chan<- lib.HostInfo) {
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
	toGatherIPsChan := make(chan lib.HostInfo)

	readerCtx, readerCancel := context.WithCancel(context.Background())

	wgGather.Add(1)

	go Gather(toGatherIPsChan, &wgGather)
	go RunWorkers(toCheckIPsChan, toGatherIPsChan)
	go ReadStdin(toCheckIPsChan, readerCtx)

	wgGather.Wait()
	readerCancel()
}
