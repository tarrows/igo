package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

const templatesDir = "templates"

type Health struct {
	Status string `json:"status"`
	Now    string `json:"Now"`
}

type Book struct {
	Title  string
	Author string
	Likes  int
}

func main() {
	var addr = flag.String("addr", ":5963", "The address of the application")
	flag.Parse()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/books", booksHandler)

	log.Println("Server listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func booksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.Path)
	books := []Book{
		{Title: "Masters of Drums", Author: "Ben Smith", Likes: 3},
		{Title: "The Smile Touch", Author: "Margaret Maximilian", Likes: 15},
		{Title: "Darkness Chain", Author: "Brandi Ni", Likes: 0},
	}
	t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "books.html")))
	_ = t.ExecuteTemplate(w, "books.html", books)
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.Path)

	health := Health{Status: "OK", Now: time.Now().String()}

	res, err := json.Marshal(health)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
