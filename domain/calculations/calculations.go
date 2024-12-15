package calculations

import (
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/wind"
)

type Calculations struct {
	CallSign         callsign.CallSign
	PressureAltitude pressure.Altitude
	OAT              temperature.Temperature
	Wind             wind.Wind

	TakeOffWeightAndBalance *WeightBalance
	LandingWeightAndBalance *WeightBalance
	FuelPlanning            *FuelPlanning

	Performance *Performance
}
