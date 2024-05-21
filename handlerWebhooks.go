package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) processingPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	params := parameter{}
	decoder := json.NewDecoder(r.Body)
	errJson := decoder.Decode(&params)
	if errJson != nil {
		respondError(w, http.StatusInternalServerError, "Couldnt decode the Parameters")
		log.Println(errJson)
		return
	}
	sentAPIKey := r.Header.Get("Authorization")

	if sentAPIKey == "" {
		respondError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}
	strippedKey, ok := strings.CutPrefix(sentAPIKey, "ApiKey ")
	if !ok {
		respondError(w, http.StatusUnauthorized, "No authorization possible with this Header")
		return
	}

	if strippedKey != cfg.APIKey {
		respondError(w, http.StatusUnauthorized, "Faulty APIKey for this request")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}
	userId := params.Data.UserId

	upgradeErr := cfg.DB.AssignMemberSatusByID(userId)
	if upgradeErr != nil {
		respondError(w, http.StatusNotFound, "No user under this ID in the database")
		return
	} else {
		respondWithJSON(w, http.StatusOK, "")
		return
	}

}
