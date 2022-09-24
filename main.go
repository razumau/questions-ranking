package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ComparisonPageData struct {
	PageTitle      string
	FirstQuestion  string
	SecondQuestion string
}

func aboutHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<html><body><h1>About this</h1></body></html>")
}

func comparisonHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/comparison.html"))
	vars := mux.Vars(r)
	data := ComparisonPageData{
		PageTitle:      vars["id"],
		FirstQuestion:  "Первый вопрос",
		SecondQuestion: "Второй вопрос",
	}
	tmpl.Execute(w, data)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/about", aboutHandler)
	r.HandleFunc("/{id}", comparisonHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
