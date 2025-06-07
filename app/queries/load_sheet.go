package queries

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/callsign"
	"github.com/michaelvlaar/ppl-calculations/domain/pressure"
	"github.com/michaelvlaar/ppl-calculations/domain/seat"
	"github.com/michaelvlaar/ppl-calculations/domain/state"
	"github.com/michaelvlaar/ppl-calculations/domain/temperature"
	"github.com/michaelvlaar/ppl-calculations/domain/weight_balance"
	"github.com/michaelvlaar/ppl-calculations/domain/wind"
)

type LoadSheetHandler struct {
	stateProvider state.Provider
}

func NewLoadSheetHandler(stateProvider state.Provider) LoadSheetHandler {
	return LoadSheetHandler{
		stateProvider: stateProvider,
	}
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

func (handler LoadSheetHandler) Handle(ctx context.Context) (LoadSheetResponse, error) {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return LoadSheetResponse{}, err
	}

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
