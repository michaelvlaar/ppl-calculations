package queries

import (
	"context"
	"errors"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/state"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/wind"
)

var (
	ErrMissingFuelSheet = errors.New("missing fuel sheet")
)

type StatsSheetHandler struct {
	calcService calculations.Service
}

func NewStatsSheetHandler(calcService calculations.Service) StatsSheetHandler {
	return StatsSheetHandler{
		calcService: calcService,
	}
}

type StatsSheetResponse struct {
	CallSign         callsign.CallSign
	PressureAltitude pressure.Altitude
	OAT              temperature.Temperature
	Wind             wind.Wind

	TakeOffWeightAndBalance *calculations.WeightBalance
	LandingWeightAndBalance *calculations.WeightBalance
	FuelPlanning            *calculations.FuelPlanning

	Performance *calculations.Performance
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

	if s.CallSign == nil {
		return sheet, ErrMissingFuelSheet
	}

	sheet.CallSign = *s.CallSign
	sheet.PressureAltitude = *s.PressureAltitude
	sheet.OAT = *s.OutsideAirTemperature
	sheet.Wind = *s.Wind

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

	_, todRR, todDR, err := handler.calcService.TakeOffDistance(*s.OutsideAirTemperature, *s.PressureAltitude, sheet.TakeOffWeightAndBalance.Total.Mass, *s.Wind, calculations.ChartTypeSVG)
	if err != nil {
		return sheet, err
	}

	_, ldrDR, ldrGR, err := handler.calcService.LandingDistance(*s.OutsideAirTemperature, *s.PressureAltitude, sheet.LandingWeightAndBalance.Total.Mass, *s.Wind, calculations.ChartTypeSVG)
	if err != nil {
		return sheet, err
	}
	sheet.Performance = &calculations.Performance{
		TakeOffRunRequired:        todRR,
		TakeOffDistanceRequired:   todDR,
		LandingDistanceRequired:   ldrDR,
		LandingGroundRollRequired: ldrGR,
	}

	return sheet, nil
}
