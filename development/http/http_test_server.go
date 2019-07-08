package main

import (
	"fmt"
	"log"
	"net/http"
)

const RedirectStatusCode = 302

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/")

	printHeader(r)

	w.Header().Add("key", "value")
}

func printHeader(r *http.Request) {
	for k, values := range r.Header {
		for _, value := range values {
			fmt.Println("key:", k, "val:", value)
		}
	}
}

func main() {
	fmt.Println("Running on localhost:9090")
	http.HandleFunc("/", httpHandler)
	http.HandleFunc("/redirect", redirectHandler)
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/redirect")
	printHeader(r)
	http.Redirect(w, r, "/", RedirectStatusCode)
}
