package queries

import (
	"context"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/state"
	"ppl-calculations/domain/volume"
	"time"
)

type FuelSheetHandler struct {
}

func NewFuelSheetHandler() FuelSheetHandler {
	return FuelSheetHandler{}
}

type FuelSheetResponse struct {
	FuelType          *fuel.Type
	FuelVolumeType    *volume.Type
	Fuel              *fuel.Fuel
	MaxFuel           *bool
	TripDuration      *time.Duration
	AlternateDuration *time.Duration
}

func (handler FuelSheetHandler) Handle(ctx context.Context, stateService state.Service) (FuelSheetResponse, error) {
	s, err := stateService.State(ctx)
	if err != nil {
		return FuelSheetResponse{}, err
	}

	return FuelSheetResponse{
		FuelType:          s.FuelType,
		FuelVolumeType:    s.FuelVolumeType,
		Fuel:              s.Fuel,
		MaxFuel:           s.MaxFuel,
		TripDuration:      s.TripDuration,
		AlternateDuration: s.AlternateDuration,
	}, nil
}
