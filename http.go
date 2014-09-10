package main

import (
	"fmt"
	"net/http"
	"strings"
	"log"
)

func HandleRootRequest(w http.ResponseWriter, r * http.Request) {
	path := r.URL.Path
	if path == "/" {
		http.ServeFile(w, r, "html/src/index.html")
	} else {
		s := []string{"html/src", path}
		path = strings.Join(s, "")
		http.ServeFile(w, r, path)
	}
}

func HandleApiRequest(w http.ResponseWriter, r * http.Request) {
	fmt.Fprint(w, "Api")
}

func httpdStart() {
	http.HandleFunc("/", HandleRootRequest)
	http.HandleFunc("/api", HandleApiRequest)
	host := "localhost:8080"
	log.Printf("Starting webserver: %s", host)
	http.ListenAndServe(host, nil)
}
