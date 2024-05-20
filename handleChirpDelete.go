package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) deleteUserChirp(w http.ResponseWriter, r *http.Request) {
	sentToken := r.Header.Get("Authorization")
	authorID, err := cfg.ValidateTokenGetId(w, sentToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid token")
		log.Println(err)
		return
	}

	chirpIDStr := r.PathValue("ID")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		log.Println("could not get the Chirp ID from the request")
		respondError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	errDelete := cfg.DB.DeleteUserChirp(authorID, chirpID)

	if errDelete != nil {
		respondError(w, http.StatusForbidden, "you do not have the rigths for this action")
	} else {
		respondWithJSON(w, http.StatusNoContent, "")
	}

}
