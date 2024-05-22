package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	database "github.com/JulianKerns/chirpy/internal/database"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get("author_id")
	sortCommand := r.URL.Query().Get("sort")

	if authorIdString == "" {
		if sortCommand == "asc" || sortCommand == "" {
			chirps, err := cfg.DB.GetChirps()
			if err != nil {
				log.Println(err)
				respondError(w, http.StatusNotFound, "this user does not exist or has no Chirps posted")
				return
			}
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })
			respondWithJSON(w, http.StatusOK, chirps)
			return

		} else if sortCommand == "desc" {
			chirps, err := cfg.DB.GetChirps()
			if err != nil {
				log.Println(err)
				respondError(w, http.StatusNotFound, "this user does not exist or has no Chirps posted")
				return
			}
			sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id > chirps[j].Id })
			respondWithJSON(w, http.StatusOK, chirps)
			return
		} else {
			respondError(w, http.StatusBadRequest, "Bad sort-query parameter")
		}
	} else {
		authorId, err := strconv.Atoi(authorIdString)
		if err != nil {
			log.Println(err)
			return
		}

		if sortCommand == "asc" || sortCommand == "" {
			userChirps, err := cfg.DB.GetUserChirps(authorId)
			if err != nil {
				log.Println(err)
				respondError(w, http.StatusNotFound, "this user does not exist or has no Chirps posted")
				return
			}
			sort.Slice(userChirps, func(i, j int) bool { return userChirps[i].Id < userChirps[j].Id })
			respondWithJSON(w, http.StatusOK, userChirps)
			return

		} else if sortCommand == "desc" {
			userChirps, err := cfg.DB.GetUserChirps(authorId)
			if err != nil {
				log.Println(err)
				respondError(w, http.StatusNotFound, "this user does not exist or has no Chirps posted")
				return
			}
			sort.Slice(userChirps, func(i, j int) bool { return userChirps[i].Id > userChirps[j].Id })
			respondWithJSON(w, http.StatusOK, userChirps)
			return
		} else {
			respondError(w, http.StatusBadRequest, "Bad sort-query parameter")
		}
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
