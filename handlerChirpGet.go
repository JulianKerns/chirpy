package main

import (
	"log"
	"net/http"
	"strconv"

	database "github.com/JulianKerns/chirpy/internal/database"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get("author_id")

	if authorIdString == "" {
		chirps, err := cfg.DB.GetChirps()
		if err != nil {
			respondError(w, http.StatusBadRequest, "Could not retrieve the Chirps from the database")
			log.Println(err)
		}
		respondWithJSON(w, http.StatusOK, chirps)
	} else {
		authorId, err := strconv.Atoi(authorIdString)
		if err != nil {
			log.Println(err)
			return
		}
		userChirps, err := cfg.DB.GetUserChirps(authorId)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusNotFound, "this user does not exist or has no Chirps posted")
		}
		respondWithJSON(w, http.StatusOK, userChirps)
	}
}

func (cfg *apiConfig) getSpecificChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("ID")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		log.Println("could not get the Chirp ID from the request")
	}

	allChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondError(w, http.StatusBadRequest, "Could not retrieve the Chirps from the database")
		log.Println(err)
	}
	emptyChirp := database.Chirp{
		Id:   0,
		Body: "",
	}
	var specificChirp database.Chirp
	for _, chirp := range allChirps {
		if chirp.Id == chirpID {
			specificChirp = chirp

		}
	}
	if specificChirp != emptyChirp {
		respondWithJSON(w, http.StatusOK, specificChirp)

	} else {

		respondError(w, http.StatusNotFound, "Chirp was not found in the database")
	}
}
