package commands

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

type UpdateLoadSheetHandler struct {
	stateProvider state.Provider
}

func NewUpdateLoadSheetHandler(stateProvider state.Provider) UpdateLoadSheetHandler {
	return UpdateLoadSheetHandler{
		stateProvider: stateProvider,
	}
}

type UpdateLoadSheetRequest struct {
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

func (handler UpdateLoadSheetHandler) Handle(ctx context.Context, request UpdateLoadSheetRequest) error {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return err
	}

	s, err := stateService.State(ctx)
	if err != nil {
		return err
	}

	s.CallSign = request.CallSign
	s.Pilot = request.Pilot
	s.PilotSeat = request.PilotSeat
	s.Passenger = request.Passenger
	s.PassengerSeat = request.PassengerSeat
	s.Baggage = request.Baggage
	s.OutsideAirTemperature = request.OutsideAirTemperature
	s.PressureAltitude = request.PressureAltitude
	s.Wind = request.Wind

	err = stateService.SetState(ctx, s)
	if err != nil {
		return err
	}

	return nil
}
