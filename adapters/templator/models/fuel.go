package models

import (
	"fmt"
	"ppl-calculations/domain/state"
)

type Fuel struct {
	Base

	FuelType          string
	FuelVolumeUnit    string
	TripDuration      *string
	AlternateDuration *string
	FuelVolume        *string
	FuelMax           bool
}

func FuelFromState(s state.State) interface{} {
	fs := Fuel{
		Base: Base{
			Step: string(StepFuel),
		},
		FuelType:       "mogas",
		FuelVolumeUnit: "liter",
		FuelMax:        false,
	}

	if s.MaxFuel != nil {
		fs.FuelMax = *s.MaxFuel
	}

	if s.FuelType != nil {
		fs.FuelType = s.FuelType.String()
	}

	if s.FuelVolumeType != nil {
		fs.FuelVolumeUnit = s.FuelVolumeType.String()
	}

	if s.Fuel != nil {
		fs.FuelVolume = StringPointer(fmt.Sprintf("%.1f", s.Fuel.Volume.Amount))
	}

	if s.TripDuration != nil {
		fs.TripDuration = StringPointer(fmt.Sprintf("%d:%d", int(s.TripDuration.Hours()), int(s.TripDuration.Minutes())%60))
	}

	if s.AlternateDuration != nil {
		fs.AlternateDuration = StringPointer(fmt.Sprintf("%d:%d", int(s.AlternateDuration.Hours()), int(s.AlternateDuration.Minutes())%60))
	}

	return fs
}
