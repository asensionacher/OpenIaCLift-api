package router

import (
	"OpenIaCLiftAPI/internal/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/login", middleware.LoginUser).Methods("POST")
	return router
}
