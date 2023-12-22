package model

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Driver struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	DriverId string  `json:"id"`
}
