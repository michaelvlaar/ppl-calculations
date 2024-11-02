package commands

import (
	"context"
	"ppl-calculations/domain/state"
)

type ClearSheetHandler struct {
}

func NewClearSheetHandler() ClearSheetHandler {
	return ClearSheetHandler{}
}

func (handler ClearSheetHandler) Handle(ctx context.Context, stateService state.Service) error {
	err := stateService.SetState(ctx, state.MustNew())
	if err != nil {
		return err
	}
	return nil
}
