package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func aboutHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "<html><body><h1>About this</h1></body></html>")
}

func ComparisonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Id: %v\n", vars["id"])
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/about", aboutHandler)
	r.HandleFunc("/{id}", ComparisonHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
