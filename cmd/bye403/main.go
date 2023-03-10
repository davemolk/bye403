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
	ignore      string
	input       bool
	insecure    bool
	method      bool
	os          string
	path        bool
	proxy       string
	redirects   bool
	rHeaders    bool
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
	ignore []int
	path   string
	sc     []int
}

func main() {
	var config config
	// techniques
	flag.BoolVar(&config.headers, "h", true, "manipulate headers")
	flag.BoolVar(&config.method, "m", true, "manipulate method")
	flag.BoolVar(&config.path, "p", true, "manipulate path")

	// config
	flag.IntVar(&config.concurrency, "c", 10, "number of concurrent requests")
	flag.BoolVar(&config.input, "i", false, "read url off stdin")
	flag.StringVar(&config.ignore, "ignore", "", "status code(s) to ignore in output (403 ignored by default)")
	flag.BoolVar(&config.insecure, "insecure", true, "accept any certificate and host name presented by server")
	flag.StringVar(&config.os, "os", "w", "operating system")
	flag.StringVar(&config.proxy, "proxy", "", "proxy to use")
	flag.BoolVar(&config.redirects, "r", true, "allow redirects")
	flag.BoolVar(&config.rHeaders, "rh", false, "include response headers in output")
	flag.BoolVar(&config.silent, "s", true, "silence error reporting")
	flag.StringVar(&config.statusCode, "sc", "", "filter output by status code(s)")
	flag.IntVar(&config.timeout, "t", 5000, "request timeout (in ms)")
	flag.StringVar(&config.url, "u", "https://www.example.com/secret", "target url")
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
		b.sc = b.statusCodes(config.statusCode)
	}
	if config.ignore != "" {
		b.ignore = b.statusCodes(config.ignore)
	}

	b.host, b.path = b.parseURL(config.url)

	b.client = b.customClient(config.proxy, config.insecure, config.redirects)

	wait := b.bye403()
	<-wait
}

func (b *bye403) bye403() <-chan struct{} {
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup
	tokens := make(chan struct{}, b.config.concurrency)
	wg.Add(3) // path, header, and method manipulation
	go func() {
		defer wg.Done()
		if b.config.path {
			paths := b.paths(b.config.url)
			for _, p := range paths {
				wg.Add(1)
				tokens <- struct{}{}
				go func(path, method string, headers []string) {
					defer wg.Done()
					url := fmt.Sprintf("https://%s%s", b.host, path)
					b.request(url, method, nil)
					<-tokens
				}(p, http.MethodGet, nil)

				wg.Add(1)
				tokens <- struct{}{}
				go func(path, method string, headers []string) {
					defer wg.Done()
					url := fmt.Sprintf("http://%s%s", b.host, path)
					b.request(url, method, headers)
					<-tokens
				}(p, http.MethodGet, nil)
			}
		}
	}()
	go func() {
		defer wg.Done()
		if b.config.headers {
			headers := b.manipulateHeaders()
			for _, h := range headers {
				wg.Add(1)
				tokens <- struct{}{}
				go func(url, method string, headers []string) {
					defer wg.Done()
					b.request(url, method, headers)
					<-tokens
				}(b.config.url, http.MethodGet, h)
			}
			wg.Add(1)
			tokens <- struct{}{}
			go func(url, method string, headers []string) {
				defer wg.Done()
				b.request(url, method, headers)
				<-tokens
			}(b.config.url, http.MethodPost, []string{"X-HTTP-Method-Override", "PUT"})
		}
	}()
	go func() {
		defer wg.Done()
		if b.config.method {
			for _, m := range b.methods() {
				wg.Add(1)
				tokens <- struct{}{}
				go func(url, method string, headers []string) {
					defer wg.Done()
					b.request(url, method, headers)
					<-tokens
				}(b.config.url, m, nil)
			}
		}
	}()
	go func() {
		defer close(done)
		wg.Wait()
	}()

	return done
}
