package main

import (
	"bufio"
	"flag"
	"fmt"
	"modular-netstalking/lib"
	"os"
)

var (
	format string
)

func init() {
	flag.StringVar(&format, "f", "%s%s", "concatenation format")
	flag.StringVar(&lib.DictPath, "d", "", "dictionary path")
}

func main() {
	flag.Parse()

	if lib.DictPath == "" {
		flag.Usage()
		os.Exit(255)
	}

	dict, err := lib.FileToArray(lib.DictPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dict open: ", err)
		os.Exit(255)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ln := scanner.Text()
		for _, item := range dict {
			fmt.Println(fmt.Sprintf(format, ln, item))
		}
	}
}
