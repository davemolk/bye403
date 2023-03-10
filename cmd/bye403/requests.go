package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func (b *bye403) request(url, method string, header []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(b.config.timeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		if !b.config.silent {
			fmt.Printf("failed to create request for %s: %v\n", url, err)
		}
		return
	}

	req = b.browserHeaders(req)

	// check if we're doing header manipulation
	var headManip bool
	if len(header) > 0 {
		headManip = true
		req.Header.Set(header[0], header[1])
	}

	// using custom client
	resp, err := b.client.Do(req)
	if err != nil {
		if !b.config.silent {
			fmt.Printf("failed to get response for %s: %v\n", url, err)
		}
		return
	}
	defer resp.Body.Close()

	for _, sc := range b.ignore {
		if resp.StatusCode == sc {
			return
		}
	}

	if resp.StatusCode != 403 {
		switch {
		case b.config.statusCode != "":
			for _, code := range b.sc {
				if code == resp.StatusCode {
					b.output(resp.StatusCode, method, header, resp.Header, url, headManip)
				}
			}
		default:
			b.output(resp.StatusCode, method, header, resp.Header, url, headManip)
		}
	}
}

func (b *bye403) output(code int, method string, reqHeaders []string, respHeaders http.Header, url string, headManip bool) {
	fmt.Printf("%d: %s request\n", code, method)
	fmt.Println(url)
	if headManip {
		fmt.Printf("headers: %s\n", reqHeaders)
	}
	if b.config.rHeaders {
		fmt.Printf("response headers: %s\n", respHeaders)
	}
	fmt.Println()
}

func (b *bye403) browserHeaders(r *http.Request) *http.Request {
	if rand.Intn(2) == 1 {
		return b.ff(r)
	}
	return b.chrome(r)
}

func (b *bye403) ff(r *http.Request) *http.Request {
	uAgent := b.ffUA()
	r.Header.Set("Host", b.host)
	r.Header.Set("User-Agent", uAgent)
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	r.Header.Set("DNT", "1")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Upgrade-Insecure-Requests", "1")
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Sec-Fetch-Site", "none")
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("Sec-GCP", "1")
	return r
}

func (b *bye403) chrome(r *http.Request) *http.Request {
	uAgent := b.chromeUA()
	r.Header.Set("Host", b.host)
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Cache-Control", "max-age=0")
	r.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"`)
	r.Header.Set("sec-ch-ua-mobile", "?0")
	switch b.config.os {
	case "m":
		r.Header.Set("sec-ch-ua-platform", "Macintosh")
	default:
		r.Header.Set("sec-ch-ua-platform", "Windows")
	}
	r.Header.Set("Upgrade-Insecure-Requests", "1")
	r.Header.Set("User-Agent", uAgent)
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	r.Header.Set("Sec-Fetch-Site", "none")
	r.Header.Set("Sec-Fetch-Mode", "navigate")
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("Sec-Fetch-Dest", "document")
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")

	return r
}

func (b *bye403) ffUA() string {
	var userAgents []string
	switch b.config.os {
	case "m":
		userAgents = []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:107.0) Gecko/20100101 Firefox/107.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:106.0) Gecko/20100101 Firefox/106.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:104.0) Gecko/20100101 Firefox/104.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:102.0) Gecko/20100101 Firefox/102.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:101.0) Gecko/20100101 Firefox/101.0",
		}
	default:
		userAgents = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0",
		}
	}
	random := rand.Intn(len(userAgents))
	return userAgents[random]
}

func (b *bye403) chromeUA() string {
	var userAgents []string
	switch b.config.os {
	case "m":
		userAgents = []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4692.56 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36",
		}
	default:
		userAgents = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
		}
	}
	random := rand.Intn(len(userAgents))
	return userAgents[random]
}
