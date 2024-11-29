package queries

import (
	"context"
	"io"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
)

type TodChartRequest struct {
	OAT              temperature.Temperature
	PressureAltitude pressure.Altitude
	Tow              weight_balance.Mass
	Wind             wind.Wind
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
	chart, todGR, todDR, err := h.calcService.TakeOffDistance(request.OAT, request.PressureAltitude, request.Tow, request.Wind)
	if err != nil {
		return nil, err
	}

	return &TodChartResponse{
		Chart:                     chart,
		TakeOffGroundRollRequired: todGR,
		TakeOffDistanceRequired:   todDR,
	}, nil
}
