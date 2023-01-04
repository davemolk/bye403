package main

import (
	"fmt"
	"net/http"
	"sync"
)

func (b *bye403)  createRequests(url string, paths []string) {
	headers := b.manipulateHeaders()

	ch := make(chan *http.Request, len(headers) * len(paths))
	var wg sync.WaitGroup
	for _, h := range headers {
		for _, p := range paths {
			wg.Add(1)
			go func (h []string, p string) {
				defer wg.Done()
				req, err := http.NewRequest(http.MethodGet, url + p, nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				req.Header.Set(h[0], h[1])
				ch <- req
			}(h, p)
		}
	}
	wg.Wait()
	close(ch)
	for c := range ch {
		fmt.Println(c)
	}
}