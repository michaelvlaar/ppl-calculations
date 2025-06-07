package queries

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/calculations"
	"github.com/michaelvlaar/ppl-calculations/domain/callsign"
	"github.com/michaelvlaar/ppl-calculations/domain/weight_balance"
	"io"
)

type WBChartHandler struct {
	service calculations.Service
}

func NewWBChartHandler(service calculations.Service) WBChartHandler {
	return WBChartHandler{
		service: service,
	}
}

type WBChartRequest struct {
	CallSign          callsign.CallSign
	TakeOffMassMoment weight_balance.MassMoment
	LandingMassMoment weight_balance.MassMoment
	WithinLimits      bool
	ChartType         calculations.ChartType
}

func (h WBChartHandler) Handle(_ context.Context, request WBChartRequest) (io.Reader, error) {
	return h.service.WeightAndBalance(request.CallSign, request.TakeOffMassMoment, request.LandingMassMoment, request.WithinLimits, request.ChartType)
}
