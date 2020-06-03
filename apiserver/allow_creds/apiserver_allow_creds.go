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
	Secret    string ",omitempty"
}

var userData = map[string]User{
	"john": {UserName: "jdoe", FirstName: "John", LastName: "Doe", Country: "France"},
}

var port = flag.Int("port", 12346, "port to listen on, default is 12346")
var allowCreds = flag.Bool("allow-cred", false, "enable to set Access-Control-Allow-Credentials")

func corsWrapper(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Printf("Request Origin header %s\n", origin)
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if *allowCreds {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

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
	user := userData[r.URL.Path[len("/users/"):]]
	if c, err := r.Cookie("token"); err == nil && c.Value == "secret_token" {
		user.Secret = "This is a user secret only shown when cookie is present"
	}
	b, _ := json.Marshal(user)
	io.WriteString(w, string(b))

}

func main() {
	flag.Parse()
	http.HandleFunc("/users/", corsWrapper(optionsWrapper(userHandler)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("apiserver.cors.com:%d", *port), nil))
}
