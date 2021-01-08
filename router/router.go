package router

import (
	"andon-datasource/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", middleware.TestConnection).Methods("GET", "OPTIONS")
	router.HandleFunc("/search", middleware.Search).Methods("POST", "OPTIONS")
	router.HandleFunc("/query", middleware.Query).Methods("POST", "OPTIONS")

	return router

}
