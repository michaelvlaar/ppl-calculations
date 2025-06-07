package commands

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/state"
)

type ClearSheetHandler struct {
	stateProvider state.Provider
}

func NewClearSheetHandler(stateProvider state.Provider) ClearSheetHandler {
	return ClearSheetHandler{
		stateProvider: stateProvider,
	}
}

func (handler ClearSheetHandler) Handle(ctx context.Context) error {
	stateService, err := handler.stateProvider.ServiceFrom(ctx)
	if err != nil {
		return err
	}

	err = stateService.SetState(ctx, state.MustNew())
	if err != nil {
		return err
	}
	return nil
}
