package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

const templatesDir = "tetmplates"

type Health struct {
	Status string `json:"status"`
	Now    string `json:"Now"`
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
