package main

import (
	"bufio"
	"flag"
	"os"
	"strings"
)

type config struct {
	timeout int
	url string
}

type bye403 struct {
	config config
}

func main() {
	var config config
	flag.StringVar(&config.url, "u", "https://www.example.com", "base url")
	flag.Parse()
	
	b := &bye403{
		config: config,
	}

	paths := b.makePaths()
	b.createRequests(config.url, paths)
	
}

func input() ([]string, error) {
	var lines []string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "q" {
			break
		}
		lines = append(lines, strings.TrimPrefix(s.Text(), "/"))
	}
	return lines, s.Err()
}