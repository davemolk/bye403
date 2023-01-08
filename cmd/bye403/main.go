package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

type config struct {
	concurrency int
	headers bool
	input bool
	method bool
	os string
	path bool
	timeout int
	url string
	validate bool
}

type bye403 struct {
	config config
	host string
	path string
}

func main() {
	var config config
	// techniques
	flag.BoolVar(&config.headers, "h", false, "manipulate headers")
	flag.BoolVar(&config.method, "m", false, "manipulate method")
	flag.BoolVar(&config.path, "p", false, "manipulate path")

	// config
	flag.IntVar(&config.concurrency, "c", 10, "number of concurrent requests")
	flag.BoolVar(&config.input, "i", false, "read url off stdin")
	flag.StringVar(&config.os, "os", "w", "operating system")
	flag.IntVar(&config.timeout, "t", 5000, "request timeout (in ms)")
	flag.StringVar(&config.url, "u", "https://www.example.com/secret", "base url")
	flag.BoolVar(&config.validate, "v", true, "validate url before running program")
	flag.Parse()
	
	b := &bye403{
		config: config,
	}

	if config.input {
		config.url = b.input()
	}

	if config.validate {
		b.validateURL(config.url)
		b.validateOS(config.os)
	}

	b.host, b.path = b.parseURL(config.url)

	var wg sync.WaitGroup
	// put in semaphore
	if config.path {
		paths := b.paths(config.url)
		for _, p := range paths {
			wg.Add(1)
			go func(path, method string, headers []string) {
				defer wg.Done()
				url := fmt.Sprintf("https://%s%s", b.host, path)
				err := b.request(url, method, nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("path success: %s\n", url)
			}(p, http.MethodGet, nil)
			wg.Add(1)
			go func(path, method string, headers []string) {
				defer wg.Done()
				url := fmt.Sprintf("http://%s%s", b.host, path)
				err := b.request(url, method, headers)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("path success: %s\n", url)
			}(p, http.MethodGet, nil)
		}
	}
	if config.headers {
		headers := b.manipulateHeaders()
		for _, h := range headers {
			wg.Add(1)
			go func(url, method string, headers[]string) {
				defer wg.Done()
				err := b.request(url, method, headers)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("headers success: %v:%v\n", headers[0], headers[1])
			}(config.url, http.MethodGet, h)
		}
	}
	if config.method {
		for _, m := range b.verbs() {
			wg.Add(1)
			go func(url, method string, headers []string) {
				defer wg.Done()
				err := b.request(url, method, headers)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("method success: %s\n", method)
			}(config.url, m, nil)
		} 
	}
	wg.Wait()
}