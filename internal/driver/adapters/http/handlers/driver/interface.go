package driverHandler

import "net/http"

type DriverHandler interface {
	ListTrips(w http.ResponseWriter, r *http.Request)
	GetTrip(w http.ResponseWriter, r *http.Request)
	CancelTrip(w http.ResponseWriter, r *http.Request)
	AcceptTrip(w http.ResponseWriter, r *http.Request)
	StartTrip(w http.ResponseWriter, r *http.Request)
	EndTrip(w http.ResponseWriter, r *http.Request)
}
