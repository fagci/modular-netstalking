package main

import (
	"bufio"
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

func CheckPort(host string, port string) bool {
	conn, _ := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

func Work(ch chan string, limited bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		ip, ok := <-ch
		if !ok {
			break
		}

		if !CheckPort(ip, port) {
			continue
		}

		printMutex.Lock()
		if count == 0 {
			printMutex.Unlock()
			break
		}

		fmt.Println(ip)

		if limited {
			count--
			if count == 0 {
				printMutex.Unlock()
				break
			}
		}
		printMutex.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	ch := make(chan string, workers*2)

	scanner := bufio.NewScanner(os.Stdin)
	limited := count > 0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go Work(ch, limited, &wg)
	}

	for scanner.Scan() {
		if count == 0 {
			break
		}
		ch <- scanner.Text()
	}

	close(ch)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	wg.Wait()
}
