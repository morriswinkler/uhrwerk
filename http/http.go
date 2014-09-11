package http

import (
	"fmt"
	"net/http"
	"log"
)

type Server struct {
	Host, Port string
}

func (s *Server) Init(Host, Port string) error {
	
	// Save config
	s.Host = Host
	s.Port = Port
	host := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// Configure handlers
	http.HandleFunc("/", s.HandleRootRequest)
	http.HandleFunc("/api", s.HandleApiRequest)
	
	// And eventually start the webserver
	log.Printf("Starting webserver: %s", host)
	err := http.ListenAndServe(host, nil)
	return err
}

// Handle Root Request
func (s *Server) HandleRootRequest(w http.ResponseWriter, r * http.Request) {
	path := r.URL.Path
	if path == "/" {
		http.ServeFile(w, r, "res/html/src/index.html")
	} else {
		path = fmt.Sprintf("res/html/src%s", path)
		http.ServeFile(w, r, path)
	}
}

func (s *Server) HandleApiRequest(w http.ResponseWriter, r * http.Request) {
	//fmt.Fprint(w, "Api")

	// Check what Negroni can do for us

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


