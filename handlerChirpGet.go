package main

import (
	"log"
	"net/http"
	"strconv"

	database "github.com/JulianKerns/chirpy/internal/database"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondError(w, http.StatusBadRequest, "Could not retrieve the Chirps from the database")
		log.Println(err)
	}
	respondWithJSON(w, http.StatusOK, chirps)

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
