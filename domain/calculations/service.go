package calculations

import (
	"io"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
)

type ChartType int

const (
	ChartTypeSVG ChartType = iota
	ChartTypePNG
)

type Service interface {
	TakeOffDistance(oat temperature.Temperature, pa pressure.Altitude, tow weight_balance.Mass, w wind.Wind, chartType ChartType) (io.Reader, float64, float64, error)
	LandingDistance(oat temperature.Temperature, pa pressure.Altitude, tow weight_balance.Mass, w wind.Wind, chartType ChartType) (io.Reader, float64, float64, error)
	WeightAndBalance(callSign callsign.CallSign, takeOffMassMoment weight_balance.MassMoment, landingMassMoment weight_balance.MassMoment, withinLimits bool, chartType ChartType) (io.Reader, error)
}
