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

	http.HandleFunc("/api/files", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListFiles(w, r)
		case http.MethodPost:
			handlers.CreateFile(w, r)
		case http.MethodDelete:
			handlers.DeleteFile(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/files/phrases", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetPhrasesFromFile(w, r)
		case http.MethodPost:
			handlers.AddPhraseToFile(w, r)
		case http.MethodDelete:
			handlers.DeletePhraseByID(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
