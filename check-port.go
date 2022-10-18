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
	count      int64
	timeout    time.Duration
	stdinMutex sync.Mutex
	countMutex sync.Mutex
)

func init() {
	flag.StringVar(&port, "p", "80", "port")
	flag.IntVar(&workers, "w", 256, "workers count")
	flag.Int64Var(&count, "c", 0, "count of IPs to gather (0 = infinite)")
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

func Work(scanner *bufio.Scanner, limited bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		stdinMutex.Lock()
		if !scanner.Scan() {
			break
		}
		ip := scanner.Text()
		stdinMutex.Unlock()
		if !CheckPort(ip, port) {
			continue
		}
		fmt.Println(ip)
		if limited {
			countMutex.Lock()
			count--
			countMutex.Unlock()
			if count <= 0 {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(255)
	}
	os.Exit(0)
}

func main() {
	var wg sync.WaitGroup
    flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	limited := count > 0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go Work(scanner, limited, &wg)
	}

	wg.Wait()
}
