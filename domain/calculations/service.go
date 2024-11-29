package calculations

import (
	"io"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
)

type Service interface {
	TakeOffDistance(oat temperature.Temperature, pa pressure.Altitude, tow weight_balance.Mass, w wind.Wind) (io.Reader, float64, float64, error)
	LandingDistance(oat temperature.Temperature, pa pressure.Altitude, tow weight_balance.Mass, w wind.Wind) (io.Reader, float64, float64, error)
}
