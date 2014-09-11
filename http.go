package main

import (
	"fmt"
	"net/http"
	"log"
)

// Handle Root Request
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
	//fmt.Fprint(w, "Api")

	// Check what Negroni can do

	// Later we will use this to accept only POST data
	fmt.Fprintf(w, "Request method: %s\n", r.Method)
	
	// We could allow only specific hosts
	fmt.Fprintf(w, "Host: %s\n", r.Host)

	// We have to do ParseForm in order to use Form as url.Values
	r.ParseForm()
	vals := r.Form
	if len(vals) > 0 {
		fmt.Fprint(w, "Variables:\n")
		for k, v := range vals {
			fmt.Fprintf(w, "%s: %s\n", k, v)
		}
	}

}

func httpdStart() {
	http.HandleFunc("/", HandleRootRequest)
	http.HandleFunc("/api", HandleApiRequest)
	host := "localhost:8080"
	log.Printf("Starting webserver: %s", host)
	http.ListenAndServe(host, nil)
}
