package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type User struct {
	UserName  string
	FirstName string
	LastName  string
	Country   string
}

var userData = map[string]User{
	"john": User{"jdoe", "John", "Doe", "France"},
}

var port = flag.Int("port", 12346, "port to listen on, default is 12346")

func corsWrapper(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Printf("Request Origin header %s\n", origin)
		w.Header().Set("Access-Control-Allow-Origin", "*") // Open all so that this is not an issue
		w.Header().Set("Access-Control-Allow-Headers", "x-tenant")
		w.Header().Set("Access-Control-Max-Age", "6000")
		if r.Method == "OPTIONS" {
			return
		}
		fn(w, r)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	uname := r.URL.Path[len("/users/"):]
	w.Header().Set("Content-Type", "application/json")
	if uname == "@me" {
		location := fmt.Sprintf("http://localhost:%d/users/john", *port)
		w.Header().Set("Location", location)
		fmt.Printf("Redirecting to %s ", location)
		w.WriteHeader(303)
		return
	} else {

		b, _ := json.Marshal(userData[uname])
		io.WriteString(w, string(b))
	}

}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "Hello World")

}

func main() {
	flag.Parse()
	http.HandleFunc("/hello", corsWrapper(helloHandler))
	http.HandleFunc("/users/", corsWrapper(userHandler))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil))
}
