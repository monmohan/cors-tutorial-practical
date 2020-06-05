package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

//var setCookie = flag.Bool("set-cookie", false, "enable to set cookie in response")
var port = flag.Int("port", 12346, "port to listen on, default is 12346")
var allowCreds = flag.Bool("allow-creds", false, "enable to set Access-Control-Allow-Credentials")
var crossSite = flag.Bool("cross-site", false, "enable to set a different site the server")
var ssModeNone = flag.Bool("same-site-none", false, "enable same site none cookie setting")

type User struct {
	UserName  string
	FirstName string
	LastName  string
	Country   string
}

var userData = map[string]User{
	"john": {UserName: "jdoe", FirstName: "John", LastName: "Doe", Country: "France"},
}

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
	dumpRequestOut(w, r)
	defer (func() { fmt.Println("-----------Request Completed----------") })()
	w.Header().Set("Content-Type", "application/json")
	uid := r.URL.Path[len("/users/"):]
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if uid == "@me" {
		//read userId from cookie
		if c, err := r.Cookie("visited-userid"); err == nil {
			fmt.Printf("Cookie detected : %v\n", *c)
			uid = c.Value
		}
	}
	if user, ok := userData[uid]; ok {

		ssMode := http.SameSiteStrictMode
		if *ssModeNone {
			ssMode = http.SameSiteNoneMode
		}
		fmt.Printf("Setting cookie for the user %v, same site value is %v\n", uid, ssMode)
		http.SetCookie(w, &http.Cookie{Name: "visited-userid", Value: uid, SameSite: ssMode})

		b, _ := json.Marshal(user)
		io.WriteString(w, string(b))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, `{"error": "user not found"}`)

}

func dumpRequestOut(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%s", dump)
}

func main() {
	flag.Parse()
	siteName := "apiserver.cors.com"
	if *crossSite {
		siteName = "apiserver.sscors.com"
	}
	http.HandleFunc("/users/", corsWrapper(optionsWrapper(userHandler)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", siteName, *port), nil))
}
