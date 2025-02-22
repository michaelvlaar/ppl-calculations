package state

import (
	"context"
)

type Provider interface {
	WithService(context.Context, Service) context.Context
	ServiceFrom(context.Context) (Service, error)
}
