package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"log"
)

var templates = template.Must(template.ParseFiles("html/dist/index.html"))
var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err

	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	fmt.Println(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func HandleRootRequest(w http.ResponseWriter, r * http.Request) {
	//fmt.Println(r.URL.Path[1:])
	//http.ServeFile(w, r, r.URL.Path[1:])

	path := r.URL.Path
	if path == "/" {
		http.ServeFile(w, r, "html/src/index.html")
	} else {
		s := []string{"html/src", path}
		path = strings.Join(s, "")
		http.ServeFile(w, r, path)
	}
}

func HandleAssetsRequest(w http.ResponseWriter, r * http.Request) {
	fmt.Println(r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
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
