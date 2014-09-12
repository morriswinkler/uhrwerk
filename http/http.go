// Fab Lab Locksmith webserver that serves the UI and provides REST API
package http

import (
	"fmt"
	"net/http"
	"net/url"
	"log"
	"strings"
	"github.com/morriswinkler/uhrwerk/debug"
	"github.com/morriswinkler/uhrwerk/database"
	"github.com/morriswinkler/uhrwerk/api"
)

// Main webserver structure
type Server struct {
	Host string
	Port string
	Dir string
	db *database.Database
	Api *api.Api
}

// Initializes the server 
func (s *Server) Init(Host, Port, Dir string, db *database.Database) error {
	
	// Save config
	s.Host = Host
	s.Port = Port
	s.Dir = Dir
	s.db = db
	
	// Init API
	s.Api = api.NewApi(s.db)

	// This is just to notify the user via log message later
	host := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// Configure handlers
	http.HandleFunc("/", s.HandleRootRequest)
	http.HandleFunc("/api/", s.HandleApiRequest)
	
	// And eventually start the webserver
	debug.INFO.Printf("Starting webserver: %s", host)
	log.Printf("Starting webserver: %s", host)
	err := http.ListenAndServe(host, nil)
	return err
}

// HandleRootRequest displays the main HTML frontend of the FabSmith
// - the Fab Lab Locksmith
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

// HandleApiRequest handles REST API cals from the outer space
func (s *Server) HandleApiRequest(w http.ResponseWriter, r * http.Request) {
	
	// Check what Negroni can do for us... Maybe... Maybe not!
	
	// Trim the url to remove leading "/" and "api"
	path := strings.TrimPrefix(r.URL.Path, "/")
	path = strings.TrimPrefix(path, "api")
	path = strings.TrimPrefix(path, "/")
	
	// We have to do ParseForm in order to use Form as url.Values
	r.ParseForm()
	var vals *url.Values = &r.Form
	method := r.Method

	// Get and serve API response
	returnString := s.Api.Call(path, method, vals)
	fmt.Fprintln(w, returnString)
}

