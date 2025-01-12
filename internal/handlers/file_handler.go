package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
			filenames = append(filenames, strings.ReplaceAll(file.Name(), ".json", ""))
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

func RenameFile(w http.ResponseWriter, r *http.Request) {
	filesMu.Lock()
	defer filesMu.Unlock()

	oldName := r.URL.Query().Get("oldName")
	newName := r.URL.Query().Get("newName")

	if oldName == "" || newName == "" {
		http.Error(w, "Both oldName and newName must be provided", http.StatusBadRequest)
		return
	}

	oldPath := fmt.Sprintf("./data/%s.json", oldName)
	newPath := fmt.Sprintf("./data/%s.json", newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		http.Error(w, "Old file does not exist", http.StatusNotFound)
		return
	}

	if _, err := os.Stat(newPath); err == nil {
		http.Error(w, "File with the new name already exists", http.StatusConflict)
		return
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		http.Error(w, "Failed to rename file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	filesMu.Lock()
	defer filesMu.Unlock()

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("./data/%s.json", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
