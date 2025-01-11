package main

import (
	"eh/internal/handlers"
	"eh/internal/utils"
	"fmt"
	"net/http"
)

var baseDir = "./data"

func main() {
	utils.CreateDataFolder(baseDir)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/phrases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetAllPhrases(w, r)
		} else if r.Method == http.MethodPost {
			handlers.AddPhrase(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/phrases/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetPhraseByID(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.DeletePhraseByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.ListFiles(w, r)
		} else if r.Method == http.MethodPost {
			handlers.CreateFile(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/files/phrases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetPhrasesFromFile(w, r)
		} else if r.Method == http.MethodPost {
			handlers.AddPhraseToFile(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
