package queries

import (
	"context"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/state"
)

type ExportsHandler struct {
	stateProvider state.Provider
}

func NewExportsHandler(stateProvider state.Provider) ExportsHandler {
	return ExportsHandler{
		stateProvider: stateProvider,
	}
}

type ExportsResponse struct {
	Exports []export.Export
}

func (handler ExportsHandler) Handle(ctx context.Context) (ExportsResponse, error) {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return ExportsResponse{}, err
	}

	resp := ExportsResponse{}

	ex, err := stateService.Exports(ctx)
	if err != nil {
		return resp, err
	}

	resp.Exports = ex

	return resp, nil
}
