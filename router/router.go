package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/windbnb/reservation-service/handler"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()

	log.Fatal(http.ListenAndServe(":8082", router))

	return router
}