package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Content string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldnt decode the Parameters")
		return
	}
	//validated, err := validateChirp(params.Content)
	//if err != nil {
	//	respondError(w, http.StatusBadRequest, "Chirpy is too long")
	//	return
	//}

	newUser, errCh := cfg.DB.CreateUser(params.Content)
	if errCh != nil {
		log.Println("could not add theUser to the database")

	}
	respondWithJSON(w, http.StatusCreated, newUser)

}
