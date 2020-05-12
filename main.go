package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic("Unable to read config.json file !")
	}

	var config Config

	jdecoder := json.NewDecoder(configFile)
	err = jdecoder.Decode(&config)
	if err != nil {
		panic("Unable to parse config.json content !")
	}

	// Setup static files serving
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve pages
	http.HandleFunc("/home/", log(IndexHandler))
	http.HandleFunc("/socket/", SocketHandler)

	if config.ENV == "production" {
		RunSecureServer(config)
	} else {
		RunServer(config)
	}
}

func RunServer(config Config) {
	server := http.Server{
		Addr: ":" + config.HTTP_PORT,
	}

	defer server.ListenAndServe()
	fmt.Printf("%s %s HTTP server started on port %s \n", config.APP_NAME, config.APP_VERSION, config.HTTP_PORT)
}

func RunSecureServer(config Config) {
	server := http.Server{
		Addr: ":" + config.HTTP_PORT,
	}

	fmt.Printf("%s %s HTTPS server started on port %s \n", config.APP_NAME, config.APP_VERSION, config.HTTP_PORT)
	defer server.ListenAndServeTLS(config.SSL_CERT_PATH, config.SSL_KEY_PATH)
}
