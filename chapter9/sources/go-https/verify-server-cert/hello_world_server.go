package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!\n")
	})
	fmt.Println(http.ListenAndServeTLS("localhost:8081",
		"/Users/wyy/code/github/https/certificate/server-signed-by-ca.crt",
		"/Users/wyy/code/github/https/certificate/server.key", nil))
}
