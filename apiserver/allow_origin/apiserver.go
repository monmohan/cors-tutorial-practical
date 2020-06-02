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
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
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
	var e error
	if e = json.NewDecoder(r.Body).Decode(&nUser); e == nil {
		fmt.Printf("User created %v\n", nUser)
		//save user
		userData[nUser.UserName] = nUser
		//write response
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(userData[nUser.UserName])
		io.WriteString(w, string(b))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Printf("Error in creating user %v", e.Error())

}

func main() {
	flag.Parse()
	http.HandleFunc("/users/", corsWrapper(userHandler))
	http.HandleFunc("/users", corsWrapper(createUserHandler))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("apiserver.cors.com:%d", *port), nil))
}
