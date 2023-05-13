package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/reservation-service/handler"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservationRequest", handler.CreateReservationRequest).Methods("POST")
	router.HandleFunc("/api/reservationRequest/guest/{id}", handler.GetGuestsActive).Methods("GET")
	router.HandleFunc("/api/reservationRequest/owner/{id}", handler.GetOwnersActive).Methods("GET")

	return router
}
