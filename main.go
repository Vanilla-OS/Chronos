package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vanilla-os/Chronos/core"
	"github.com/vanilla-os/Chronos/settings"
)

var version = "0.2.0"

func main() {
	err := core.LoadChronos()
	if err != nil {
		log.Printf("unable to load Chronos: %s", err.Error())
		os.Exit(1)
	}

	// CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			next.ServeHTTP(w, r)
		})
	}

	// Define router with routes and their handlers
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "version": "` + version + `"}`))
	})
	r.HandleFunc("/repos", core.HandleRepos)
	r.HandleFunc("/{repoId}", core.HandleRepo)
	r.HandleFunc("/{repoId}/langs", core.HandleLangs)
	r.HandleFunc("/{repoId}/articles/{lang}", core.HandleArticles)
	r.HandleFunc("/{repoId}/articles/{lang}/{slug}", core.HandleArticle)
	r.HandleFunc("/{repoId}/search/{lang}", core.HandleSearch)

	http.Handle("/", r)

	// Start the server
	log.Printf("Chronos listening on port %s...\n", settings.Cnf.Port)
	log.Printf("Address: http://0.0.0.0:%s\n", settings.Cnf.Port)
	log.Println("Press Ctrl+C to exit.")
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", settings.Cnf.Port), nil)
}
