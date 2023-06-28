package router

import (
	"github.com/gorilla/mux"
	"github.com/windbnb/reservation-service/handler"
	"github.com/windbnb/reservation-service/metrics"
)

func ConfigureRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/reservationRequest/new", metrics.MetricProxy(handler.CreateReservationRequest)).Methods("POST")
	router.HandleFunc("/api/reservationRequest/guest/{id}", metrics.MetricProxy(handler.GetGuestsActive)).Methods("GET")
	router.HandleFunc("/api/reservationRequest/owner/{id}", metrics.MetricProxy(handler.GetOwnersActive)).Methods("GET")
	router.HandleFunc("/api/reservationRequest/{id}", metrics.MetricProxy(handler.DeleteReservationRequest)).Methods("DELETE")
	router.HandleFunc("/api/reservationRequest/{id}/accept", metrics.MetricProxy(handler.AcceptReservationRequest)).Methods("PUT")
	router.HandleFunc("/api/reservationRequest/{id}/cancel", metrics.MetricProxy(handler.CancelReservationRequest)).Methods("PUT")
	router.HandleFunc("/api/reservationRequest/{guestId}/cancelled", metrics.MetricProxy(handler.CountGuestsCancelledReservations)).Methods("GET")
	router.HandleFunc("/api/reservationRequest/guest/{id}/all", metrics.MetricProxy(handler.GetGuestsReservations)).Methods("GET")
	router.HandleFunc("/api/reservationRequest/owners/{id}", metrics.MetricProxy(handler.GetOwnersReservations)).Methods("GET")

	router.HandleFunc("/api/reservationRequest/guest/{guestId}/host/{hostId}", handler.GetWheatherGuestWasWithHost).Methods("GET")
	router.HandleFunc("/api/reservationRequest/guest/{guestId}/accomodation/{accomodationId}", handler.GetWheatherGuestWasInAccomodation).Methods("GET")

	router.HandleFunc("/probe/liveness", handler.Healthcheck)
	router.HandleFunc("/probe/readiness", handler.Ready)

	router.Path("/metrics").Handler(metrics.MetricsHandler())

	return router
}
