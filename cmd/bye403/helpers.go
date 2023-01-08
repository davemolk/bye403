package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
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

func (b *bye403) input() string {
	s := bufio.NewScanner(os.Stdin) 
	var u string
	for s.Scan() {
		u = s.Text()
	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}
	return u
}

func (b *bye403) statusCodes() []int {
	if b.config.statusCode == "" {
		return nil
	}
	codes := strings.Split(b.config.statusCode, " ")
	sc := make([]int, len(codes))
	for _, code := range codes {
		c, err := strconv.Atoi(code)
		if err != nil {
			log.Fatal(err)
		}
		sc = append(sc, c)
	}
	return sc
}