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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding the request parameters %s", err)
		w.WriteHeader(500)
		return
	}
	clean := cleanBody(params.Content)
	contentLength := len(clean)
	log.Println(contentLength)
	type responseVal struct {
		Valid string `json:"cleaned_body"`
	}
	type responseErr struct {
		Error string `json:"error"`
	}

	if contentLength > 140 {
		answerBody := responseErr{
			Error: "Chirp is too long",
		}

		data, err := json.Marshal(answerBody)
		if err != nil {
			log.Printf("Error decoding the request parameters %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(data))

	} else {
		answerBody := responseVal{
			Valid: clean,
		}
		data, err := json.Marshal(answerBody)
		if err != nil {
			log.Printf("Error decoding the request parameters %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(data))
	}

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
