package queries

import (
	"context"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/state"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
)

type LoadSheetHandler struct {
}

func NewLoadSheetHandler() LoadSheetHandler {
	return LoadSheetHandler{}
}

type LoadSheetResponse struct {
	CallSign              *callsign.CallSign
	Pilot                 *weight_balance.Mass
	PilotSeat             *seat.Position
	Passenger             *weight_balance.Mass
	PassengerSeat         *seat.Position
	Baggage               *weight_balance.Mass
	OutsideAirTemperature *temperature.Temperature
	PressureAltitude      *pressure.Altitude
	Wind                  *wind.Wind
}

func (handler LoadSheetHandler) Handle(ctx context.Context, stateService state.Service) (LoadSheetResponse, error) {
	s, err := stateService.State(ctx)
	if err != nil {
		return LoadSheetResponse{}, err
	}

	return LoadSheetResponse{
		CallSign:              s.CallSign,
		Pilot:                 s.Pilot,
		PilotSeat:             s.PilotSeat,
		Passenger:             s.Passenger,
		PassengerSeat:         s.PassengerSeat,
		Baggage:               s.Baggage,
		OutsideAirTemperature: s.OutsideAirTemperature,
		PressureAltitude:      s.PressureAltitude,
		Wind:                  s.Wind,
	}, nil
}
