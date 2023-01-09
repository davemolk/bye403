package main

import "net/http"

// skip MethodDelete
func (b *bye403) verbs() []string {
	return []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
}
