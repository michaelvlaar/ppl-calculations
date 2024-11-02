package commands

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

type UpdateLoadSheetHandler struct {
}

func NewUpdateLoadSheetHandler() UpdateLoadSheetHandler {
	return UpdateLoadSheetHandler{}
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

func (handler UpdateLoadSheetHandler) Handle(ctx context.Context, stateService state.Service, request UpdateLoadSheetRequest) error {
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
