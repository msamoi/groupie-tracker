package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"groupie-tracker/data"
	"groupie-tracker/unmarshal"
)

func main() {
	mux := http.NewServeMux()
	server := http.FileServer(http.Dir("./static"))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/artist/", artist)
	mux.Handle("/static/", http.StripPrefix("/static/", server))
	fmt.Printf("Starting server at port 8080, press CTRL + c to terminate the process.\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found.", http.StatusNotFound)
		return
	}
	page, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "Unable to reach the server - 500 Internal Server Error.", http.StatusInternalServerError)
		return
	}
	unmarshal.Fetch(data.APIL, 0)
	err = page.Execute(w, &data.Artists)
	if err != nil {
		http.Error(w, "Unable to reach the server - 500 Internal Server Error.", http.StatusInternalServerError)
		return
	}
}

func artist(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`^(/artist/\d+)$`)
	if !re.MatchString(r.URL.Path) {
		http.Error(w, "404 Page Not Found.", http.StatusNotFound)
		return
	}
	page, err := template.ParseFiles("static/artist.html")
	if err != nil {
		http.Error(w, "Unable to reach the server - 500 Internal Server Error.", http.StatusInternalServerError)
		return
	}

	pageid := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, pagerr := strconv.Atoi(pageid)
	if pagerr != nil {
		http.Error(w, "404 Page Not Found.", http.StatusNotFound)
		return
	}
	unmarshal.Fetch(data.APIL, id)
	err = page.Execute(w, data.Artists[id-1])
	if err != nil {
		http.Error(w, "Unable to reach the server - 500 Internal Server Error.", http.StatusInternalServerError)
		return
	}
}
