package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/reservation-service/handler"
	"github.com/windbnb/reservation-service/metrics"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservationRequest", handler.CreateReservationRequest).Methods("POST")
	router.HandleFunc("/api/reservationRequest/guest/{id}", handler.GetGuestsActive).Methods("GET")
	router.HandleFunc("/api/reservationRequest/owner/{id}", handler.GetOwnersActive).Methods("GET")
	router.HandleFunc("/api/reservationRequest/{id}", handler.DeleteReservationRequest).Methods("DELETE")
	router.HandleFunc("/api/reservationRequest/{id}", handler.AcceptReservationRequest).Methods("PUT")
	router.HandleFunc("/api/reservationRequest/{id}/cancel", handler.CancelReservationRequest).Methods("PUT")
	router.HandleFunc("/api/reservationRequest/{guestId}/cancelled", handler.CountGuestsCancelledReservations).Methods("GET")

	router.HandleFunc("/api/reservationRequest/guest/{guestId}/host/{hostId}", handler.GetWheatherGuestWasWithHost).Methods("GET")
	router.HandleFunc("/api/reservationRequest/guest/{guestId}/accomodation/{accomodationId}", handler.GetWheatherGuestWasInAccomodation).Methods("GET")

	router.Path("/metrics").Handler(metrics.MetricsHandler())

	return router
}
