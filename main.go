package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	database "github.com/JulianKerns/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	const filepath string = "/home/julian_k/workspace/github.com/JulianKerns/GoProjects/chirpy/database.json"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		os.Remove(filepath)
	}

	databaseJSON, err := database.NewDB(filepath)
	if err != nil {
		log.Fatalln("could not create the database.json file")
	}
	config := apiConfig{
		fileServerHits: 0,
		DB:             databaseJSON,
	}

	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app")))))

	mux.Handle("GET /api/metrics", config.displayMetrics())

	mux.Handle("/api/reset", config.resetMetrics())

	mux.HandleFunc("GET /admin/metrics", config.adminMetrics)

	mux.HandleFunc("GET /api/healthz", readinessHandler)

	mux.HandleFunc("POST /api/chirps", config.postChirp)

	mux.HandleFunc("GET /api/chirps", config.getChirps)

	mux.HandleFunc("GET /api/chirps/{ID...}", config.getSpecificChirp)

	mux.HandleFunc("POST /api/users", config.createUser)

	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
