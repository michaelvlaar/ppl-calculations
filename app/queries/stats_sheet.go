package queries

import (
	"context"
	"errors"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/state"
)

var (
	ErrMissingFuelSheet = errors.New("missing fuel sheet")
)

type StatsSheetHandler struct {
}

func NewStatsSheetHandler() StatsSheetHandler {
	return StatsSheetHandler{}
}

type StatsSheetResponse struct {
	CallSign                callsign.CallSign
	TakeOffWeightAndBalance *calculations.WeightBalance
	LandingWeightAndBalance *calculations.WeightBalance
	FuelPlanning            *calculations.FuelPlanning
}

func (handler StatsSheetHandler) Handle(ctx context.Context, stateService state.Service) (StatsSheetResponse, error) {
	sheet := StatsSheetResponse{}

	s, err := stateService.State(ctx)
	if err != nil {
		return sheet, err
	}

	if s.TripDuration == nil {
		return sheet, ErrMissingFuelSheet
	}

	sheet.CallSign = *s.CallSign

	var f fuel.Fuel
	if s.MaxFuel != nil && *s.MaxFuel {
		sheet.TakeOffWeightAndBalance, f, err = calculations.NewWeightAndBalanceMaxFuel(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, *s.Baggage, *s.FuelType)
		if err != nil {
			return sheet, err
		}
	} else {
		f = *s.Fuel
		sheet.TakeOffWeightAndBalance, err = calculations.NewWeightAndBalance(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, *s.Baggage, f)
		if err != nil {
			return sheet, err
		}
	}

	sheet.FuelPlanning, err = calculations.NewFuelPlanning(*s.TripDuration, *s.AlternateDuration, f, *s.FuelVolumeType)
	if err != nil {
		return sheet, err
	}

	sheet.LandingWeightAndBalance, err = calculations.NewWeightAndBalance(*s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, *s.Baggage, fuel.Subtract(f, sheet.FuelPlanning.Trip, sheet.FuelPlanning.Taxi))
	if err != nil {
		return sheet, err
	}

	return sheet, nil
}
