package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	comparisonTemplate = template.Must(template.ParseFiles("templates/comparison.html"))
	db                 = DB()
)

type ComparisonPageData struct {
	ComparisonId   string
	FirstQuestion  string
	SecondQuestion string
}

func DB() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://localhost:5432/questions")
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}

func createComparison() string {
	var id string
	query := "insert into comparisons (first, second) values (1, 2) returning id"
	err := db.QueryRow(context.Background(), query).Scan(&id)
	if err != nil {
		log.Println("Could not create a comparison")
	}
	return id

}

func saveComparison(id string, selection string) {
	var query string
	log.Println(selection)
	switch selection {
	case "first":
		query = "update comparisons set selected = first where id = $1"
	case "second":
		query = "update comparisons set selected = second where id = $1"
	default:
		log.Println("wrong value for selection parameter")
	}

	res, err := db.Exec(context.Background(), query, id)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

func fetchComparison(id string) ComparisonPageData {
	db.QueryRow(context.Background(), `select c.id, q1.question, q2.question 
	from comparisons c
	join questions q1 on q1.id = c.first
	join questions q2 on q2.id = c.second
	where c.id = $1`, id)
	return ComparisonPageData{
		ComparisonId:   id,
		FirstQuestion:  "Первый вопрос",
		SecondQuestion: "Второй вопрос",
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><body><h1>About this</h1></body></html>")
}

func selectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comparison_id := vars["id"]
	selection := vars["selection"]

	saveComparison(comparison_id, selection)

	next_comparison := createComparison()
	http.Redirect(w, r, "/"+next_comparison, http.StatusFound)
}

func comparisonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := fetchComparison(vars["id"])
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

	defer db.Close()

	log.Fatal(srv.ListenAndServe())
}
