package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Web World!")
}

func main() {
	http.HandleFunc("/", handler1)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
