package commands

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/export"
	"github.com/michaelvlaar/ppl-calculations/domain/state"
)

type DeleteExportSheetHandler struct {
	stateProvider state.Provider
}

func NewDeleteExportSheetHandler(stateProvider state.Provider) DeleteExportSheetHandler {
	return DeleteExportSheetHandler{
		stateProvider: stateProvider,
	}
}

type DeleteExportSheetRequest struct {
	ID export.ID
}

func (handler DeleteExportSheetHandler) Handle(ctx context.Context, request DeleteExportSheetRequest) error {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return err
	}

	return stateService.DeleteExport(ctx, request.ID)
}
