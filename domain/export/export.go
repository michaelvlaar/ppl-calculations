package export

import (
	"github.com/michaelvlaar/ppl-calculations/domain/callsign"
	"github.com/michaelvlaar/ppl-calculations/domain/fuel"
	"github.com/michaelvlaar/ppl-calculations/domain/pressure"
	"github.com/michaelvlaar/ppl-calculations/domain/seat"
	"github.com/michaelvlaar/ppl-calculations/domain/temperature"
	"github.com/michaelvlaar/ppl-calculations/domain/weight_balance"
	"github.com/michaelvlaar/ppl-calculations/domain/wind"
	"time"
)

type Export struct {
	ID   ID
	Name Name

	CallSign              callsign.CallSign
	Pilot                 weight_balance.Mass
	PilotSeat             seat.Position
	Passenger             *weight_balance.Mass
	PassengerSeat         *seat.Position
	Baggage               *weight_balance.Mass
	OutsideAirTemperature temperature.Temperature
	PressureAltitude      pressure.Altitude
	Wind                  wind.Wind

	Fuel fuel.Fuel

	TripDuration      time.Duration
	AlternateDuration time.Duration

	CreatedAt time.Time
}

func New(id ID, name Name, callSign callsign.CallSign, pilot weight_balance.Mass, pilotSeat seat.Position, passenger *weight_balance.Mass, passengerSeat *seat.Position, baggage *weight_balance.Mass, outsideAirTemperature temperature.Temperature, pressureAltitude pressure.Altitude, wind wind.Wind, fuel fuel.Fuel, tripDuration time.Duration, alternateDuration time.Duration, createdAt time.Time) (Export, error) {
	return Export{
		ID:                    id,
		Name:                  name,
		CallSign:              callSign,
		Pilot:                 pilot,
		PilotSeat:             pilotSeat,
		Passenger:             passenger,
		PassengerSeat:         passengerSeat,
		Baggage:               baggage,
		OutsideAirTemperature: outsideAirTemperature,
		PressureAltitude:      pressureAltitude,
		Wind:                  wind,
		Fuel:                  fuel,
		TripDuration:          tripDuration,
		AlternateDuration:     alternateDuration,
		CreatedAt:             createdAt,
	}, nil
}
