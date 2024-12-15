package queries

import (
	"context"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
)

type ExportSheetHandler struct {
}

func NewExportSheetHandler() ExportSheetHandler {
	return ExportSheetHandler{}
}

type ExportSheetResponse struct {
	Name *export.Name
}

func (handler ExportSheetHandler) Handle(ctx context.Context, stateService state.Service) (ExportSheetResponse, error) {
	sheet := ExportSheetResponse{}

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
