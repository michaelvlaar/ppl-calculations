package commands

import (
	"context"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
)

type DeleteExportSheetHandler struct {
}

func NewDeleteExportSheetHandler() DeleteExportSheetHandler {
	return DeleteExportSheetHandler{}
}

type DeleteExportSheetRequest struct {
	ID export.ID
}

func (handler DeleteExportSheetHandler) Handle(ctx context.Context, stateService state.Service, request DeleteExportSheetRequest) error {
	return stateService.DeleteExport(ctx, request.ID)
}
