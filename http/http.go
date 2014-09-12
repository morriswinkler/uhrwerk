// Fab Lab Locksmith webserver that serves the UI and provides REST API
package http

import (
	"fmt"
	"net/http"
	"log"
)

// Main webserver structure
type Server struct {
	Host, Port, Dir string
}

// Initializes the server 
func (s *Server) Init(Host, Port, Dir string) error {
	
	// Save config
	s.Host = Host
	s.Port = Port
	s.Dir = Dir
	host := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// Configure handlers
	http.HandleFunc("/", s.HandleRootRequest)
	http.HandleFunc("/api", s.HandleApiRequest)
	
	// And eventually start the webserver
	log.Printf("Starting webserver: %s", host)
	err := http.ListenAndServe(host, nil)
	return err
}

// Handle Root Request - show Fab Lab Locksmith main UI
func (s *Server) HandleRootRequest(w http.ResponseWriter, r * http.Request) {
	path := r.URL.Path
	if path == "/" {
		fileToServe := fmt.Sprintf("%s/index.html", s.Dir)
		http.ServeFile(w, r, fileToServe)
	} else {
		path = fmt.Sprintf("%s%s", s.Dir, path)
		http.ServeFile(w, r, path)
	}
}

// Api call handler
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


