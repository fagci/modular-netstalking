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
	port    string
	workers int
	timeout time.Duration
)

func init() {
	flag.StringVar(&port, "p", "80", "port")
	flag.IntVar(&workers, "w", 256, "workers count")
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

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mutex.Lock()
				if !scanner.Scan() {
					break
				}
				ip := scanner.Text()
				mutex.Unlock()
				if CheckPort(ip, port) {
					fmt.Println(ip)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				os.Exit(255)
			}
			os.Exit(0)
		}()
	}

	wg.Wait()
}
