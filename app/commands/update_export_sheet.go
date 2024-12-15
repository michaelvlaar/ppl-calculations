package commands

import (
	"context"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
	"time"
)

type UpdateExportSheetHandler struct {
}

func NewUpdateExportSheetHandler() UpdateExportSheetHandler {
	return UpdateExportSheetHandler{}
}

type UpdateExportSheetRequest struct {
	ID   export.ID
	Name export.Name
}

func (handler UpdateExportSheetHandler) Handle(ctx context.Context, stateService state.Service, request UpdateExportSheetRequest) error {
	s, err := stateService.State(ctx)
	if err != nil {
		return err
	}

	e, err := export.New(request.ID, request.Name, *s.CallSign, *s.Pilot, *s.PilotSeat, s.Passenger, s.PassengerSeat, s.Baggage, *s.OutsideAirTemperature, *s.PressureAltitude, *s.Wind, *s.Fuel, *s.TripDuration, *s.AlternateDuration, time.Now())
	if err != nil {
		return err
	}

	err = stateService.SetExport(ctx, e)
	if err != nil {
		return err
	}

	return nil
}
