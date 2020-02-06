package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const templatesDir = "templates"
const assetsDir = "assets"

type Health struct {
	Status string `json:"status"`
	Now    string `json:"Now"`
}

type Book struct {
	ID     int
	Title  string
	Author string
	Likes  int
}

var books = []Book{
	{ID: 4321345, Title: "Masters of Drums", Author: "Ben Smith", Likes: 3},
	{ID: 6678453, Title: "The Smile Touch", Author: "Margaret Maximilian", Likes: 15},
	{ID: 3245561, Title: "Darkness Chain", Author: "Brandi Ni", Likes: 0},
}

func main() {
	var addr = flag.String("addr", ":5963", "The address of the application")
	flag.Parse()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/books/", bookItemHandler)
	http.HandleFunc("/books", bookListHandler)
	http.HandleFunc("/book", redirectHandler)

	log.Println("Server listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.Path)
	http.Redirect(w, req, "/books", http.StatusFound)
}

func bookItemHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.Path)

	//       segs[0]   [1]   [2]
	// localhost:3215/books/:id
	segs := strings.Split(req.URL.Path, "/")

	if len(segs[2]) < 1 {
		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "book.list.html")))
		_ = t.ExecuteTemplate(w, "book.list.html", books)
		return
	}

	id, err := strconv.Atoi(segs[2])
	if err != nil {
		log.Println("Parse Failed:", err)
		http.NotFound(w, req)
		return
	}

	for _, book := range books {
		if book.ID == id {
			t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "book.item.html")))
			_ = t.ExecuteTemplate(w, "book.item.html", book)
			return
		}
	}

	http.NotFound(w, req)
}

func bookListHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.Path)

	t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "book.list.html")))
	_ = t.ExecuteTemplate(w, "book.list.html", books)
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
