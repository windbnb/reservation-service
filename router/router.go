package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/windbnb/reservation-service/handler"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservationRequest", handler.CreateReservationRequest).Methods("POST")

	log.Fatal(http.ListenAndServe(":8083", router))

	return router
}