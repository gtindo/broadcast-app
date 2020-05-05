package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := "4000"

	// Setup static files serving
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve pages
	http.HandleFunc("/home/", log(IndexHandler))
	http.HandleFunc("/socket/", SocketHandler)

	RunServer(port)
}

func RunServer(port string) {
	server := http.Server{
		Addr: ":" + port,
	}

	defer server.ListenAndServe()
	fmt.Printf("HTTP server started on port %s \n", port)
}
