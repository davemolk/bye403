package main

import (
	"log"
	"net/url"
)

func (b *bye403) validateURL(target string) {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}
	if u.Host == "" && u.Scheme == "" {
		log.Fatal("input URL does not appear to be valid")
	}
}

func (b *bye403) validateOS(os string) {
	if os != "m" && os != "w" {
		log.Fatal("input OS does not appear to be valid")
	}
}

func (b *bye403) parseURL(target string) (string, string) {
	u, err := url.ParseRequestURI(target)
	if err != nil {
		log.Fatal(err)
	}
	return u.Host, u.Path
}