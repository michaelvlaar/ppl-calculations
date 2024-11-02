package state

import "context"

type Service interface {
	SetState(context.Context, *State) error
	State(ctx context.Context) (*State, error)
}
