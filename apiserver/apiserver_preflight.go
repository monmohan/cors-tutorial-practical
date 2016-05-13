package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
		log.Printf("Request Origin header %s\n", origin)
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		fn(w, r)
	}
}

func optionsWrapper(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqMethod, reqHeader := r.Header.Get("Access-Control-Request-Method"), r.Header.Get("Access-Control-Request-Headers")
		//check for validity
		if (r.Method == "OPTIONS") && (reqMethod == "GET" || reqMethod == "POST") && (strings.EqualFold(reqHeader, "Content-Type")) {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}

		fn(w, r)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(userData[r.URL.Path[len("/users/"):]])
	io.WriteString(w, string(b))

}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	//post only
	if r.Method != "POST" {
		return
	}
	var nUser User
	if e := json.NewDecoder(r.Body).Decode(&nUser); e == nil {
		log.Printf("User created %v", nUser)
		userData[nUser.UserName] = nUser
	}

}

func main() {
	flag.Parse()
	http.HandleFunc("/users/", corsWrapper(userHandler))
	http.HandleFunc("/users", corsWrapper(optionsWrapper(createUserHandler)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil))
}
