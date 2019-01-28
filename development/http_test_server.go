package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Println("key:", k, "val:", v)
	}

	w.Header().Add("key", "value")
	w.Header().Add("key", "value")
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}