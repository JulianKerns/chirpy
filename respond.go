package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func validateChirp(body string) (string, error) {
	const characterMax int = 140
	if len(body) > characterMax {
		return "", errors.New("chirp is too long")

	}
	clean := cleanBody(body)
	return clean, nil

}
func respondError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type responseErr struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, responseErr{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error decoding the request parameters %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func cleanBody(msg string) string {
	words := strings.Fields(msg)

	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
