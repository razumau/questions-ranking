package main

import (
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var comparisonTemplate = template.Must(template.ParseFiles("templates/comparison.html"))

type ComparisonPageData struct {
	ComparisonId   string
	FirstQuestion  string
	SecondQuestion string
}

func buildNextId() string {
	return "/15"
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><body><h1>About this</h1></body></html>")
}

func selectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comparison_id := vars["id"]
	selection := vars["selection"]
	log.Printf("Saving that for %s the %s question has been selected", comparison_id, selection)
	http.Redirect(w, r, buildNextId(), http.StatusFound)
}

func comparisonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := ComparisonPageData{
		ComparisonId:   vars["id"],
		FirstQuestion:  "Первый вопрос",
		SecondQuestion: "Второй вопрос",
	}
	comparisonTemplate.Execute(w, data)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/about", aboutHandler)
	r.Path("/submission/{id}").Queries("selection", "{selection}").HandlerFunc(selectionHandler)
	r.HandleFunc("/{id}", comparisonHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
