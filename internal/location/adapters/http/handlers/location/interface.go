package locationHandler

import "net/http"

type LocationHandler interface {
	UpdateLocation(w http.ResponseWriter, r *http.Request)
	GetNearbyDrivers(w http.ResponseWriter, r *http.Request)
}
