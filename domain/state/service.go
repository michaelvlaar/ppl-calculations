package state

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/export"
)

type Service interface {
	SetState(context.Context, *State) error
	State(context.Context) (*State, error)

	SetExport(context.Context, export.Export) error
	DeleteExport(context.Context, export.ID) error
	Export(context.Context, export.ID) (*export.Export, error)
	Exports(context.Context) ([]export.Export, error)
}
