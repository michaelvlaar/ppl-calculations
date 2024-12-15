package calculations

import (
	"io"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
	"time"
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
	Calculations(callSign *callsign.CallSign, pilot *weight_balance.Mass, pilotSeat *seat.Position, passenger *weight_balance.Mass, passengerSeat *seat.Position, baggage *weight_balance.Mass, outsideAirTemperature *temperature.Temperature, pa *pressure.Altitude, w *wind.Wind, f *fuel.Fuel, tripDuration *time.Duration, alternateDuration *time.Duration) (*Calculations, error)
}
