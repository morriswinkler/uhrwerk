package main

import (
	"fmt"
	"net/http"
	"log"
)

func HandleRootRequest(w http.ResponseWriter, r * http.Request) {
	path := r.URL.Path
	if path == "/" {
		http.ServeFile(w, r, "html/src/index.html")
	} else {
		path = fmt.Sprintf("html/src%s", path)
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
