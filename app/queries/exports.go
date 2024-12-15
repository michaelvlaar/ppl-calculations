package queries

import (
	"context"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
)

type ExportsHandler struct {
}

func NewExportsHandler() ExportsHandler {
	return ExportsHandler{}
}

type ExportsResponse struct {
	Exports []export.Export
}

func (handler ExportsHandler) Handle(ctx context.Context, stateService state.Service) (ExportsResponse, error) {
	resp := ExportsResponse{}

	ex, err := stateService.Exports(ctx)
	if err != nil {
		return resp, err
	}

	resp.Exports = ex

	return resp, nil
}
