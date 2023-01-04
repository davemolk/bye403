package main

import (
	"log"
	"net/url"
)

func (b *bye403) validate(target string) {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}
	if u.Host == "" && u.Scheme == "" {
		log.Fatal("input URL does not appear to be valid")
	}
}