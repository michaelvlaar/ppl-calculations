package state

import (
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/volume"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
	"time"
)

type State struct {
	// Load Sheet
	CallSign              *callsign.CallSign
	Pilot                 *weight_balance.Mass
	PilotSeat             *seat.Position
	Passenger             *weight_balance.Mass
	PassengerSeat         *seat.Position
	Baggage               *weight_balance.Mass
	OutsideAirTemperature *temperature.Temperature
	PressureAltitude      *pressure.Altitude
	Wind                  *wind.Wind

	// Fuel Sheet
	FuelType          *fuel.Type
	FuelVolumeType    *volume.Type
	Fuel              *fuel.Fuel
	MaxFuel           *bool
	TripDuration      *time.Duration
	AlternateDuration *time.Duration

	// Export
	ExportName *export.Name
}

func MustNew() *State {
	return &State{}
}
