package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("port", 10001, "help message for flagname")

func fileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Requested URL %v\n", r.URL.Path)
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "secret_token"})
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	flag.Parse()
	http.HandleFunc("/", fileHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil))
}
