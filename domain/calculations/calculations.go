package calculations

import (
	"github.com/michaelvlaar/ppl-calculations/domain/callsign"
	"github.com/michaelvlaar/ppl-calculations/domain/pressure"
	"github.com/michaelvlaar/ppl-calculations/domain/temperature"
	"github.com/michaelvlaar/ppl-calculations/domain/wind"
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
