package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondError(w, http.StatusBadRequest, "Could not retrieve the Chirps from the database")
		log.Println(err)
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
