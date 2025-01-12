package handlers

import (
	"eh/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	nextID = 1
)

func GetPhrasesFromFile(w http.ResponseWriter, r *http.Request) {
	filesMu.Lock()
	defer filesMu.Unlock()

	fileName := r.URL.Query().Get("file") + ".json"
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
	fileName := r.URL.Query().Get("file") + ".json"
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	var phrase models.Phrase
	if err := json.NewDecoder(r.Body).Decode(&phrase); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	phrase.ID = nextID
	nextID++

	filePath := filepath.Join(baseDir, fileName)
	filesMu.Lock()
	defer filesMu.Unlock()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	var phrases []models.Phrase
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

func DeletePhraseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	fileName := r.URL.Query().Get("file") + ".json"
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("./data/%s", fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	var phrases []models.Phrase
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&phrases); err != nil {
		http.Error(w, "Failed to parse file", http.StatusInternalServerError)
		return
	}

	phraseFound := false
	for i, p := range phrases {
		if p.ID == id {
			phrases = append(phrases[:i], phrases[i+1:]...)
			phraseFound = true
			break
		}
	}

	if !phraseFound {
		http.Error(w, "Phrase not found", http.StatusNotFound)
		return
	}

	file, err = os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(phrases); err != nil {
		http.Error(w, "Failed to write to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
