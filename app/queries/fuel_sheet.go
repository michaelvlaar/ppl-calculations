package queries

import (
	"context"
	"errors"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/state"
	"ppl-calculations/domain/volume"
	"time"
)

var (
	ErrMissingLoadSheet = errors.New("missing load sheet")
)

type FuelSheetHandler struct {
	stateProvider state.Provider
}

func NewFuelSheetHandler(stateProvider state.Provider) FuelSheetHandler {
	return FuelSheetHandler{
		stateProvider: stateProvider,
	}
}

type FuelSheetResponse struct {
	FuelType          *fuel.Type
	FuelVolumeType    *volume.Type
	Fuel              *fuel.Fuel
	MaxFuel           *bool
	TripDuration      *time.Duration
	AlternateDuration *time.Duration
}

func (handler FuelSheetHandler) Handle(ctx context.Context) (FuelSheetResponse, error) {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return FuelSheetResponse{}, err
	}

	s, err := stateService.State(ctx)
	if err != nil {
		return FuelSheetResponse{}, err
	}

	if s.CallSign == nil {
		return FuelSheetResponse{}, ErrMissingLoadSheet
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
