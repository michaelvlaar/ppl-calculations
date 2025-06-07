package calculations

import (
	"github.com/michaelvlaar/ppl-calculations/domain/callsign"
	"github.com/michaelvlaar/ppl-calculations/domain/fuel"
	"github.com/michaelvlaar/ppl-calculations/domain/pressure"
	"github.com/michaelvlaar/ppl-calculations/domain/seat"
	"github.com/michaelvlaar/ppl-calculations/domain/temperature"
	"github.com/michaelvlaar/ppl-calculations/domain/weight_balance"
	"github.com/michaelvlaar/ppl-calculations/domain/wind"
	"io"
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
