package state

import (
	"encoding/base64"
	"encoding/json"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/volume"
	"ppl-calculations/domain/weight"
	"ppl-calculations/domain/wind"
	"time"
)

type State struct {
	// Weight
	CallSign              *callsign.CallSign       `json:"callSign,omitempty"`
	Pilot                 *weight.Weight           `json:"pilot,omitempty"`
	PilotSeat             *seat.Position           `json:"pilotSeat,omitempty"`
	Passenger             *weight.Weight           `json:"passenger,omitempty"`
	PassengerSeat         *seat.Position           `json:"passengerSeat,omitempty"`
	Baggage               *weight.Weight           `json:"baggage,omitempty"`
	OutsideAirTemperature *temperature.Temperature `json:"outsideAirTemperature,omitempty"`
	PressureAltitude      *pressure.Altitude       `json:"pressureAltitude,omitempty"`
	Wind                  *wind.Wind               `json:"wind,omitempty"`

	// Fuel
	FuelType          *fuel.Type     `json:"fuelType,omitempty"`
	FuelVolumeType    *volume.Type   `json:"fuelVolumeType,omitempty"`
	Fuel              *fuel.Fuel     `json:"fuel,omitempty"`
	MaxFuel           *bool          `json:"maxFuel,omitempty"`
	TripDuration      *time.Duration `json:"tripDuration,omitempty"`
	AlternateDuration *time.Duration `json:"alternateDuration,omitempty"`
}

func (state *State) String() string {
	jsonState, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(jsonState)
}

func MustNew() *State {
	return &State{}
}

func NewFromString(state string) (*State, error) {
	s := &State{}

	base64DecodedState, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(base64DecodedState, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
