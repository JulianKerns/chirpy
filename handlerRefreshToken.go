package main

import (
	"log"
	"net/http"
	"strings"

	database "github.com/JulianKerns/chirpy/internal/database"
)

func (cfg *apiConfig) refreshingAccess(w http.ResponseWriter, r *http.Request) {
	sentToken := r.Header.Get("Authorization")
	strippedRefreshToken, ok := strings.CutPrefix(sentToken, "Bearer ")
	if !ok {
		log.Println("could not remove prefix")
	}
	var specificUser database.DatabaseUser
	var match bool

	allUsers, err := cfg.DB.GetUsers()
	if err != nil {
		log.Println(err)
		return
	}

	for _, user := range allUsers {
		if strippedRefreshToken == user.RefreshToken {
			specificUser = user
			match = true
		}
	}

	if match {
		accessToken, err := cfg.createToken(0, specificUser.Id)
		if err != nil {
			log.Println(err)
			return
		}
		respondWithJSON(w, http.StatusOK, database.RespondUser{Token: accessToken})

	} else {
		respondError(w, http.StatusUnauthorized, "no matching RefreshToken found")
	}

}
