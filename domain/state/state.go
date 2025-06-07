package state

import (
	"github.com/michaelvlaar/ppl-calculations/domain/callsign"
	"github.com/michaelvlaar/ppl-calculations/domain/export"
	"github.com/michaelvlaar/ppl-calculations/domain/fuel"
	"github.com/michaelvlaar/ppl-calculations/domain/pressure"
	"github.com/michaelvlaar/ppl-calculations/domain/seat"
	"github.com/michaelvlaar/ppl-calculations/domain/temperature"
	"github.com/michaelvlaar/ppl-calculations/domain/volume"
	"github.com/michaelvlaar/ppl-calculations/domain/weight_balance"
	"github.com/michaelvlaar/ppl-calculations/domain/wind"
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
