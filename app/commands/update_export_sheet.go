package commands

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/export"
	"github.com/michaelvlaar/ppl-calculations/domain/state"
	"time"
)

type UpdateExportSheetHandler struct {
	stateProvider state.Provider
}

func NewUpdateExportSheetHandler(stateProvider state.Provider) UpdateExportSheetHandler {
	return UpdateExportSheetHandler{
		stateProvider: stateProvider,
	}
}

type UpdateExportSheetRequest struct {
	ID   export.ID
	Name export.Name
}

func (handler UpdateExportSheetHandler) Handle(ctx context.Context, request UpdateExportSheetRequest) error {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return err
	}

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
