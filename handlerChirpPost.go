package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Content string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldnt decode the Parameters")
		return
	}
	sentToken := r.Header.Get("Authorization")

	validatedUserId, err := cfg.ValidateTokenGetId(w, sentToken)
	if err != nil {
		log.Println(err)
		return
	}

	validated, err := validateChirp(params.Content)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Chirpy is too long")
		return
	}

	newChirp, errCh := cfg.DB.CreateUserChirp(validated, validatedUserId)

	if errCh != nil {
		log.Println("could not add the Chirp to the database")

	}
	respondWithJSON(w, http.StatusCreated, newChirp)

}
