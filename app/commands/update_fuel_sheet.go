package commands

import (
	"context"
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

	err = stateService.SetState(ctx, s)
	if err != nil {
		return err
	}

	return nil
}
