package queries

import (
	"context"
	"github.com/michaelvlaar/ppl-calculations/domain/calculations"
	"github.com/michaelvlaar/ppl-calculations/domain/pressure"
	"github.com/michaelvlaar/ppl-calculations/domain/temperature"
	"github.com/michaelvlaar/ppl-calculations/domain/weight_balance"
	"github.com/michaelvlaar/ppl-calculations/domain/wind"
	"io"
)

type TodChartRequest struct {
	OAT              temperature.Temperature
	PressureAltitude pressure.Altitude
	Tow              weight_balance.Mass
	Wind             wind.Wind
	ChartType        calculations.ChartType
}

type TodChartHandler struct {
	calcService calculations.Service
}

func NewTodChartHandler(calcService calculations.Service) TodChartHandler {
	return TodChartHandler{
		calcService: calcService,
	}
}

type TodChartResponse struct {
	Chart                     io.Reader
	TakeOffDistanceRequired   float64
	TakeOffGroundRollRequired float64
}

func (h TodChartHandler) Handle(_ context.Context, request TodChartRequest) (*TodChartResponse, error) {
	chart, todGR, todDR, err := h.calcService.TakeOffDistance(request.OAT, request.PressureAltitude, request.Tow, request.Wind, request.ChartType)
	if err != nil {
		return nil, err
	}

	return &TodChartResponse{
		Chart:                     chart,
		TakeOffGroundRollRequired: todGR,
		TakeOffDistanceRequired:   todDR,
	}, nil
}
