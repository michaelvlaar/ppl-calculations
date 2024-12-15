package commands

import (
	"context"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/state"
	"ppl-calculations/domain/volume"
	"time"
)

type UpdateFuelSheetHandler struct {
}

func NewUpdateFuelSheetHandler() UpdateFuelSheetHandler {
	return UpdateFuelSheetHandler{}
}

type UpdateFuelSheetRequest struct {
	FuelType          *fuel.Type
	FuelVolumeType    *volume.Type
	Fuel              *fuel.Fuel
	MaxFuel           *bool
	TripDuration      *time.Duration
	AlternateDuration *time.Duration
}

func (handler UpdateFuelSheetHandler) Handle(ctx context.Context, stateService state.Service, request UpdateFuelSheetRequest) error {
	s, err := stateService.State(ctx)
	if err != nil {
		return err
	}

	s.FuelType = request.FuelType
	s.FuelVolumeType = request.FuelVolumeType
	s.Fuel = request.Fuel
	s.MaxFuel = request.MaxFuel
	s.TripDuration = request.TripDuration
	s.AlternateDuration = request.AlternateDuration

	var f fuel.Fuel
	if s.MaxFuel != nil && *s.MaxFuel {
		_, f, err = calculations.NewWeightAndBalanceMaxFuel(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, *s.Baggage, *s.FuelType)
		if err != nil {
			return err
		}
	} else {
		f = *s.Fuel
		_, err = calculations.NewWeightAndBalance(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, s.Baggage, f)
		if err != nil {
			return err
		}
	}

	s.Fuel = &f

	err = stateService.SetState(ctx, s)
	if err != nil {
		return err
	}

	return nil
}
