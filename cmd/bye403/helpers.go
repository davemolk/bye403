package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net/http"
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

func (b *bye403) customClient(proxy string, insecure, redirect bool) *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	switch {
	case proxy != "":
		parsed, err := url.Parse(proxy)
		if err != nil {
			log.Fatal(err)
		}

		tr.Proxy = http.ProxyURL(parsed)

		return &http.Client{
			CheckRedirect: b.allowRedirects(redirect),
			Transport:     tr,
		}
	default:
		return &http.Client{
			CheckRedirect: b.allowRedirects(redirect),
			Transport:     tr,
		}
	}
}

func (b *bye403) allowRedirects(redirect bool) func(*http.Request, []*http.Request) error {
	switch {
	case redirect:
		return nil
	default:
		return func(req *http.Request, via []*http.Request) error {
			log.Printf("blocked attempted redirect to %s\n", req.URL.String())
			return http.ErrUseLastResponse
		}
	}
}
