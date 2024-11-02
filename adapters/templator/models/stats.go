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
		template.FuelTaxi = planning.Taxi.Volume.String(*s.FuelVolumeType)
		template.FuelTrip = planning.Trip.Volume.String(*s.FuelVolumeType)
		template.FuelAlternate = planning.Alternate.Volume.String(*s.FuelVolumeType)
		template.FuelContingency = planning.Contingency.Volume.String(*s.FuelVolumeType)
		template.FuelReserve = planning.Reserve.Volume.String(*s.FuelVolumeType)
		template.FuelTotal = planning.Total.Volume.String(*s.FuelVolumeType)
		template.FuelExtra = planning.Extra.Volume.String(*s.FuelVolumeType)
		template.FuelExtraAbs = planning.Extra.Volume.String(*s.FuelVolumeType)

	} else {
		panic("not implemented")
	}

	return template
}
