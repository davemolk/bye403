package main

import (
	"flag"
	"fmt"
)

type config struct {
	headers bool
	method bool
	path bool
	timeout int
	url string
	validate bool
}

type bye403 struct {
	config config
}

func main() {
	var config config
	flag.BoolVar(&config.headers, "h", false, "manipulate headers")
	flag.BoolVar(&config.method, "m", false, "manipulate method")
	flag.BoolVar(&config.path, "p", false, "manipulate path")
	flag.IntVar(&config.timeout, "t", 5000, "request timeout (in ms)")
	flag.StringVar(&config.url, "u", "https://www.example.com", "base url")
	flag.BoolVar(&config.validate, "v", true, "validate url before running program")
	flag.Parse()
	
	b := &bye403{
		config: config,
	}

	if config.validate {
		b.validate(config.url)
	}

	paths := b.paths(config.url)
	for _, p := range paths {
		fmt.Println(p)
	}
	// b.createRequests(config.url, paths)
	
}