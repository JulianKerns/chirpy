package main

import (
	"net/http"
)

func (cfg *apiConfig) resetMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		cfg.fileServerHits = 0
		w.Write([]byte("Hits reset to 0"))
	})
}
