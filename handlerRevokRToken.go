package main

import (
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) revokeUserRToken(w http.ResponseWriter, r *http.Request) {

	sentRToken := r.Header.Get("Authorization")
	strippedRToken, ok := strings.CutPrefix(sentRToken, "Bearer ")
	if !ok {
		log.Println("could not remove prefix")
	}
	errRevoke := cfg.DB.RevokeRToken(strippedRToken)

	if errRevoke != nil {
		respondError(w, http.StatusUnauthorized, "could not revoke the RefreshToken")
	} else {
		respondWithJSON(w, http.StatusNoContent, "")
	}

}
