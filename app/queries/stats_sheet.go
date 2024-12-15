package queries

import (
	"context"
	"errors"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
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
	Calculations *calculations.Calculations
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

	c, err := handler.calcService.Calculations(s.CallSign, s.Pilot, s.PilotSeat, s.Passenger, s.PassengerSeat, s.Baggage, s.OutsideAirTemperature, s.PressureAltitude, s.Wind, s.Fuel, s.TripDuration, s.AlternateDuration)
	if err != nil {
		return sheet, err
	}

	sheet.Calculations = c

	return sheet, err
}

func (handler StatsSheetHandler) HandleExport(_ context.Context, e export.Export) (StatsSheetResponse, error) {
	sheet := StatsSheetResponse{}

	c, err := handler.calcService.Calculations(&e.CallSign, &e.Pilot, &e.PilotSeat, e.Passenger, e.PassengerSeat, e.Baggage, &e.OutsideAirTemperature, &e.PressureAltitude, &e.Wind, &e.Fuel, &e.TripDuration, &e.AlternateDuration)
	if err != nil {
		return sheet, err
	}

	sheet.Calculations = c

	return sheet, err
}
