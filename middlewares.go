package main

import (
	"fmt"
	"net/http"
)

func log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		h(w, r)
	}
}
