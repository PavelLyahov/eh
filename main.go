package main

import (
	"eh/internal/handlers"
	"eh/internal/utils"
	"embed"
	"fmt"
	"net/http"
)

//go:embed static/*
var content embed.FS

var baseDir = "./data"

func main() {
	utils.CreateDataFolder(baseDir)
	includeStaticFiles()

	http.HandleFunc("/api/files", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListFiles(w, r)
		case http.MethodPost:
			handlers.CreateFile(w, r)
		case http.MethodDelete:
			handlers.DeleteFile(w, r)
		case http.MethodPatch:
			handlers.RenameFile(w, r)
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

func includeStaticFiles() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			data, _ := content.ReadFile("static/index.html")
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
		} else {
			data, err := content.ReadFile("static" + r.URL.Path)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			switch {
			case r.URL.Path[len(r.URL.Path)-4:] == ".css":
				w.Header().Set("Content-Type", "text/css")
			case r.URL.Path[len(r.URL.Path)-3:] == ".js":
				w.Header().Set("Content-Type", "application/javascript")
			}
			w.Write(data)
		}
	})
}
