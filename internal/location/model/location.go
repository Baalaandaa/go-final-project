package model

type LatLngLiteral struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius float64 `json:"radius"`
}

type Driver struct {
	Lat      float64 `json:"lat" db:"lat"`
	Lng      float64 `json:"lng" db:"lng"`
	DriverId string  `json:"id"  db:"driver_id"`
}
