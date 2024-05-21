package main

import (
	"flag"

	"log"
	"net/http"
	"os"

	database "github.com/JulianKerns/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	JWTSecret      string
}

func main() {
	err := godotenv.Load("secret.env")
	if err != nil {
		log.Println("could not load the environment variables")
	}
	jwtSecret := os.Getenv("JWT_SECRET")

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
		JWTSecret:      jwtSecret,
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

	mux.HandleFunc("POST /api/refresh", config.refreshingAccess)

	mux.HandleFunc("POST /api/revoke", config.revokeUserRToken)

	mux.HandleFunc("DELETE /api/chirps/{ID...}", config.deleteUserChirp)

	mux.HandleFunc("POST /api/login", config.loginUserToken)

	mux.HandleFunc("PUT /api/users", config.updateUser)

	mux.HandleFunc("POST /api/polka/webhooks", config.processingPolkaWebhook)

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
