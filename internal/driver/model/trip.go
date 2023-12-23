package model

type LatLngLiteral struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

type Money struct {
	Amount   float64 `json:"amount" bson:"amount"`
	Currency string  `json:"currency" bson:"currency"`
}

type Trip struct {
	Id           string        `json:"id" bson:"id"`
	DriverId     string        `json:"driver_id" bson:"driver_id"`
	From         LatLngLiteral `json:"from" bson:"from"`
	To           LatLngLiteral `json:"to" bson:"to"`
	Price        Money         `json:"price" bson:"price"`
	Status       string        `json:"status" bson:"status"`
	CancelReason *string       `json:"cancel_reason,omitempty" bson:"cancel_reason,omitempty"`
}
