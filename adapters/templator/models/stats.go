package models

import (
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/state"
)

type Stats struct {
	Base

	FuelSufficient  bool
	FuelTaxi        string
	FuelTrip        string
	FuelAlternate   string
	FuelContingency string
	FuelReserve     string
	FuelTotal       string
	FuelExtra       string
	FuelExtraAbs    string
}

func StatsFromState(s state.State) interface{} {
	template := Stats{
		Base: Base{
			Step: string(StepStats),
		},
	}

	if s.MaxFuel != nil && *s.MaxFuel {
		planning, err := calculations.NewMaxFuelPlanning(*s.FuelType, *s.TripDuration, *s.AlternateDuration)
		if err != nil {
			panic(err)
		}

		template.FuelSufficient = planning.Sufficient
		template.FuelTaxi = planning.Taxi.Volume.String()
		template.FuelTrip = planning.Trip.Volume.String()
		template.FuelAlternate = planning.Alternate.Volume.String()
		template.FuelContingency = planning.Contingency.Volume.String()
		template.FuelReserve = planning.Reserve.Volume.String()
		template.FuelTotal = planning.Total.Volume.String()
		template.FuelExtra = planning.Extra.Volume.String()
		template.FuelExtraAbs = planning.Extra.Volume.String()

	} else {
		panic("not implemented")
	}

	return template
}
