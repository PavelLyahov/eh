package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Phrase struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
	Level       int    `json:"level"`
}

var (
	phrases   []Phrase
	phrasesMu sync.Mutex
	nextID    = 1
)

func GetAllPhrases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	phrasesMu.Lock()
	defer phrasesMu.Unlock()
	json.NewEncoder(w).Encode(phrases)
}

func AddPhrase(w http.ResponseWriter, r *http.Request) {
	var p Phrase
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	phrasesMu.Lock()
	defer phrasesMu.Unlock()

	p.ID = nextID
	nextID++
	phrases = append(phrases, p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func GetPhraseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	phrasesMu.Lock()
	defer phrasesMu.Unlock()

	for _, p := range phrases {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.NotFound(w, r)
}

func DeletePhraseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	phrasesMu.Lock()
	defer phrasesMu.Unlock()

	for i, p := range phrases {
		if p.ID == id {
			phrases = append(phrases[:i], phrases[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.NotFound(w, r)
}
