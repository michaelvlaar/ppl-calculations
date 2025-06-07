package queries

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/export"
	"github.com/michaelvlaar/ppl-calculations/domain/state"
)

type ExportSheetHandler struct {
	stateProvider state.Provider
}

func NewExportSheetHandler(stateProvider state.Provider) ExportSheetHandler {
	return ExportSheetHandler{
		stateProvider: stateProvider,
	}
}

type ExportSheetResponse struct {
	Name *export.Name
}

func (handler ExportSheetHandler) Handle(ctx context.Context) (ExportSheetResponse, error) {
	sheet := ExportSheetResponse{}

	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return sheet, err
	}

	s, err := stateService.State(ctx)
	if err != nil {
		return sheet, err
	}

	if s.MaxFuel == nil {
		return sheet, ErrMissingFuelSheet
	}

	sheet.Name = s.ExportName

	return sheet, nil
}
