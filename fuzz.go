package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	path   string
	format string
)

func init() {
	flag.StringVar(&path, "d", "", "dictionary path")
	flag.StringVar(&format, "f", "%s%s", "concatenation format")
}

func FileToArray(path string) ([]string, error) {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func main() {
	flag.Parse()

	if path == "" {
		flag.Usage()
		os.Exit(255)
	}

	dict, err := FileToArray(path)
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
