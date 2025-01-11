package models

type Phrase struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
	Level       int    `json:"level"`
}
