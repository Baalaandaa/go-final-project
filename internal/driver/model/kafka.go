package model

import (
	kafka_producer "final-project/pkg/kafka-producer"
	"github.com/gofrs/uuid"
	"time"
)

type CreatedTripEvent struct {
	ID              string `json:"id"`
	Source          string `json:"source"`
	Type            string `json:"type"`
	DataContentType string `json:"datacontenttype"`
	Time            string `json:"time"`
	Data            *Trip  `json:"data"`
}

func (ev CreatedTripEvent) ToKafkaMessage() *kafka_producer.KafkaMessage {
	return &kafka_producer.KafkaMessage{
		ID:              ev.ID,
		Source:          ev.Source,
		Type:            ev.Type,
		DataContentType: ev.DataContentType,
		Time:            ev.Time,
		Data:            ev.Data,
	}
}

func NewCreatedTripEvent(data *Trip) *CreatedTripEvent {
	id, _ := uuid.NewV4()
	return &CreatedTripEvent{
		ID:              id.String(),
		Source:          "/driver",
		Type:            "trip.event.create",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data:            data,
	}
}

type AcceptTripPayload struct {
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
}

type AcceptTripCommand struct {
	ID              string             `json:"id"`
	Source          string             `json:"source"`
	Type            string             `json:"type"`
	DataContentType string             `json:"datacontenttype"`
	Time            string             `json:"time"`
	Data            *AcceptTripPayload `json:"data"`
}

func (cmd AcceptTripCommand) ToKafkaMessage() *kafka_producer.KafkaMessage {
	return &kafka_producer.KafkaMessage{
		ID:              cmd.ID,
		Source:          cmd.Source,
		Type:            cmd.Type,
		DataContentType: cmd.DataContentType,
		Time:            cmd.Time,
		Data:            cmd.Data,
	}
}

func NewAcceptTripCommand(tripID, driverID string) *AcceptTripCommand {
	id, _ := uuid.NewV4()
	return &AcceptTripCommand{
		ID:              id.String(),
		Source:          "/driver",
		Type:            "trip.command.accept",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: &AcceptTripPayload{
			TripID:   tripID,
			DriverID: driverID,
		},
	}
}

type EndTripPayload struct {
	TripID string `json:"trip_id"`
}

type EndTripCommand struct {
	ID              string          `json:"id"`
	Source          string          `json:"source"`
	Type            string          `json:"type"`
	DataContentType string          `json:"datacontenttype"`
	Time            string          `json:"time"`
	Data            *EndTripPayload `json:"data"`
}

func (cmd EndTripCommand) ToKafkaMessage() *kafka_producer.KafkaMessage {
	return &kafka_producer.KafkaMessage{
		ID:              cmd.ID,
		Source:          cmd.Source,
		Type:            cmd.Type,
		DataContentType: cmd.DataContentType,
		Time:            cmd.Time,
		Data:            cmd.Data,
	}
}

func NewEndTripCommand(tripID string) *EndTripCommand {
	id, _ := uuid.NewV4()
	return &EndTripCommand{
		ID:              id.String(),
		Source:          "/driver",
		Type:            "trip.command.end",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: &EndTripPayload{
			TripID: tripID,
		},
	}
}

type StartTripPayload struct {
	TripID string `json:"trip_id"`
}

type StartTripCommand struct {
	ID              string            `json:"id"`
	Source          string            `json:"source"`
	Type            string            `json:"type"`
	DataContentType string            `json:"datacontenttype"`
	Time            string            `json:"time"`
	Data            *StartTripPayload `json:"data"`
}

func (cmd StartTripCommand) ToKafkaMessage() *kafka_producer.KafkaMessage {
	return &kafka_producer.KafkaMessage{
		ID:              cmd.ID,
		Source:          cmd.Source,
		Type:            cmd.Type,
		DataContentType: cmd.DataContentType,
		Time:            cmd.Time,
		Data:            cmd.Data,
	}
}

func NewStartTripCommand(tripID string) *StartTripCommand {
	id, _ := uuid.NewV4()
	return &StartTripCommand{
		ID:              id.String(),
		Source:          "/driver",
		Type:            "trip.command.end",
		DataContentType: "application/json",
		Time:            time.Now().Format(time.RFC3339),
		Data: &StartTripPayload{
			TripID: tripID,
		},
	}
}
