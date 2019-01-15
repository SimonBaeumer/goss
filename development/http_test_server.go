package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
)

const PORT = 8080

func handler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Fatal(e)
	}

	log.Println("Body", string(body))
	log.Println("Header", r.Header)

	w.Header().Add("key", "value")
}

func main() {
	log.Println("Starting http server on port ", PORT)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}