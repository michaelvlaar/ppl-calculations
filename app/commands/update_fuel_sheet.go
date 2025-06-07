package commands

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/calculations"
	"github.com/michaelvlaar/ppl-calculations/domain/fuel"
	"github.com/michaelvlaar/ppl-calculations/domain/state"
	"github.com/michaelvlaar/ppl-calculations/domain/volume"
	"time"
)

type UpdateFuelSheetHandler struct {
	stateProvider state.Provider
}

func NewUpdateFuelSheetHandler(stateProvider state.Provider) UpdateFuelSheetHandler {
	return UpdateFuelSheetHandler{
		stateProvider: stateProvider,
	}
}

type UpdateFuelSheetRequest struct {
	FuelType          *fuel.Type
	FuelVolumeType    *volume.Type
	Fuel              *fuel.Fuel
	MaxFuel           *bool
	TripDuration      *time.Duration
	AlternateDuration *time.Duration
}

func (handler UpdateFuelSheetHandler) Handle(ctx context.Context, request UpdateFuelSheetRequest) error {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return err
	}

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
