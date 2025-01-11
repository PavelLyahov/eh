package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var filesMu sync.Mutex
var baseDir = "./data"

func ListFiles(w http.ResponseWriter, r *http.Request) {
	filesMu.Lock()
	defer filesMu.Unlock()

	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			filenames = append(filenames, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filenames)
}

func CreateFile(w http.ResponseWriter, r *http.Request) {
	filesMu.Lock()
	defer filesMu.Unlock()

	var req struct {
		FileName string `json:"fileName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FileName == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(baseDir, req.FileName+".json")
	if _, err := os.Stat(filePath); err == nil {
		http.Error(w, "File already exists", http.StatusConflict)
		return
	}

	if err := ioutil.WriteFile(filePath, []byte("[]"), 0644); err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetPhrasesFromFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(baseDir, fileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func AddPhraseToFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	var phrase Phrase
	if err := json.NewDecoder(r.Body).Decode(&phrase); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(baseDir, fileName)
	filesMu.Lock()
	defer filesMu.Unlock()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var phrases []Phrase
	if err := json.Unmarshal(data, &phrases); err != nil {
		http.Error(w, "Failed to parse file", http.StatusInternalServerError)
		return
	}

	phrases = append(phrases, phrase)

	newData, err := json.Marshal(phrases)
	if err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	if err := ioutil.WriteFile(filePath, newData, 0644); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
