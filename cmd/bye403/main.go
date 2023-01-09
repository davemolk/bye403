package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

type config struct {
	concurrency int
	headers     bool
	input       bool
	insecure    bool
	method      bool
	os          string
	path        bool
	proxy       string
	redirects   bool
	silent      bool
	statusCode  string
	timeout     int
	url         string
	validate    bool
}

type bye403 struct {
	client *http.Client
	config config
	host   string
	path   string
	sc     []int
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
	flag.BoolVar(&config.insecure, "insecure", true, "accept any certificate and host name presented by server")
	flag.StringVar(&config.os, "os", "w", "operating system")
	flag.StringVar(&config.proxy, "proxy", "", "proxy to use")
	flag.BoolVar(&config.redirects, "r", true, "allow redirects")
	flag.BoolVar(&config.silent, "s", true, "silence error reporting")
	flag.StringVar(&config.statusCode, "sc", "", "filter output by status code")
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

	if config.statusCode != "" {
		b.sc = b.statusCodes()
	}

	b.host, b.path = b.parseURL(config.url)

	b.client = b.customClient(config.proxy, config.insecure, config.redirects)

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
				}
			}(p, http.MethodGet, nil)
			wg.Add(1)
			go func(path, method string, headers []string) {
				defer wg.Done()
				url := fmt.Sprintf("http://%s%s", b.host, path)
				err := b.request(url, method, headers)
				if err != nil {
					fmt.Println(err)
				}
			}(p, http.MethodGet, nil)
		}
	}
	if config.headers {
		headers := b.manipulateHeaders()
		for _, h := range headers {
			wg.Add(1)
			go func(url, method string, headers []string) {
				defer wg.Done()
				err := b.request(url, method, headers)
				if err != nil {
					fmt.Println(err)
				}
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
				}
			}(config.url, m, nil)
		}
	}
	wg.Wait()
}
