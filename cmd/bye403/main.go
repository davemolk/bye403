package main

import (
	"flag"
	"fmt"
)

type config struct {
	headers bool
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
	flag.BoolVar(&config.headers, "h", false, "manipulate headers")
	flag.BoolVar(&config.method, "m", false, "manipulate method")
	flag.StringVar(&config.os, "os", "w", "operating system")
	flag.BoolVar(&config.path, "p", false, "manipulate path")
	flag.IntVar(&config.timeout, "t", 5000, "request timeout (in ms)")
	flag.StringVar(&config.url, "u", "https://www.example.com/secret", "base url")
	flag.BoolVar(&config.validate, "v", true, "validate url before running program")
	flag.Parse()
	
	b := &bye403{
		config: config,
	}

	if config.validate {
		b.validateURL(config.url)
		b.validateOS(config.os)
	}

	b.host, b.path = b.parseURL(config.url)

	if config.path {
		paths := b.paths(config.url)
		for _, p := range paths {
			fmt.Println(p)
		}
	}
	if config.headers {
		headers := b.manipulateHeaders()
		for _, h := range headers {
			fmt.Println(h)
		}
	}

}