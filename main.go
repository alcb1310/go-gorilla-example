package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/0x4149/logz"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello World")
}
func AuthRoutes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Auth Routes")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logz.Info(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logz.Info("USER AUTH: ", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		// next.ServeHTTP(w, r)
		w.WriteHeader(http.StatusUnauthorized)
		return
	})
}

func init() {
	logz.VerbosMode()

	logz.Run()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("", HomeHandler)
	// adding a site wide middleware
	r.Use(loggingMiddleware)

	// adding a route speciific middleware
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("", AuthRoutes)
	auth.Use(authMiddleware)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "42069"
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
	}

	logz.Info("Server running on port ", port)
	logz.Fatal(srv.ListenAndServe())
}
