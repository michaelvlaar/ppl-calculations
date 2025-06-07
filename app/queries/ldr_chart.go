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

type LdrChartRequest struct {
	OAT              temperature.Temperature
	PressureAltitude pressure.Altitude
	Tow              weight_balance.Mass
	Wind             wind.Wind
	ChartType        calculations.ChartType
}

type LdrChartResponse struct {
	Chart                     io.Reader
	LandingGroundRollRequired float64
	LandingDistanceRequired   float64
}
type LdrChartHandler struct {
	calcService calculations.Service
}

func NewLdrChartHandler(calcService calculations.Service) LdrChartHandler {
	return LdrChartHandler{
		calcService: calcService,
	}
}

func (h LdrChartHandler) Handle(_ context.Context, request LdrChartRequest) (*LdrChartResponse, error) {
	chart, ldrDR, ldrGR, err := h.calcService.LandingDistance(request.OAT, request.PressureAltitude, request.Tow, request.Wind, request.ChartType)
	if err != nil {
		return nil, err
	}

	return &LdrChartResponse{
		Chart:                     chart,
		LandingDistanceRequired:   ldrDR,
		LandingGroundRollRequired: ldrGR,
	}, nil
}
