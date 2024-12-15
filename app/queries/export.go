package queries

import (
	"context"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
)

type ExportHandler struct {
}

type ExportHandlerRequest struct {
	ID export.ID
}

func NewExportHandler() ExportHandler {
	return ExportHandler{}
}

func (handler ExportHandler) Handle(ctx context.Context, stateService state.Service, request ExportHandlerRequest) (*export.Export, error) {
	return stateService.Export(ctx, request.ID)
}
