package main

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanilla-os/Chronos/core"
	"github.com/vanilla-os/Chronos/settings"
)

var version = "0.2.0"

func main() {
	core.LoadChronos()

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
	r.HandleFunc("/{repoId}", core.HandleRepo)
	r.HandleFunc("/{repoId}/langs", core.HandleLangs)
	r.HandleFunc("/{repoId}/articles/{lang}", core.HandleArticles)
	r.HandleFunc("/{repoId}/articles/{lang}/{slug}", core.HandleArticle)
	r.HandleFunc("/{repoId}/search/{lang}", core.HandleSearch)

	http.Handle("/", r)

	// Start the server
	fmt.Printf("Server listening on port %s...\n", settings.Cnf.Port)
	fmt.Printf("Address: http://localhost:%s\n", settings.Cnf.Port)
	fmt.Println("Press Ctrl+C to exit.")
	http.ListenAndServe(fmt.Sprintf(":%s", settings.Cnf.Port), nil)
}
