package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Content string `json:"body"`
	}
	const characterMax int = 140

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldnt decode the Parameters")
		return
	}
	clean := cleanBody(params.Content)
	contentLength := len(clean)

	type responseVal struct {
		Valid string `json:"cleaned_body"`
	}

	if contentLength > characterMax {
		respondError(w, http.StatusBadRequest, "Chirpy is too long")
		return

	}
	respondWithJSON(w, http.StatusOK, responseVal{
		Valid: clean})
}

func respondError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Respongind with 5XX error: %s", msg)
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
