package models

import (
	"fmt"
	"net/http"
	"ppl-calculations/app/queries"
)

type Fuel struct {
	FuelType          string
	FuelVolumeUnit    string
	TripDuration      *string
	AlternateDuration *string
	FuelVolume        *string
	FuelMax           bool
}

type FuelOption struct {
	FuelVolumeUnit string
	FuelVolume     *string
	FuelMax        bool
}

func FuelOptionFromRequest(r *http.Request) FuelOption {
	fs := FuelOption{}

	fs.FuelMax = r.URL.Query().Get("fuel_max") == "max"
	fs.FuelVolumeUnit = r.URL.Query().Get("fuel_unit")
	fs.FuelVolume = StringPointer(r.URL.Query().Get("fuel_volume"))

	return fs
}

func FuelFromFuelSheet(s queries.FuelSheetResponse) Fuel {
	fs := Fuel{
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
		fs.TripDuration = StringPointer(fmt.Sprintf("%02d%02d", int(s.TripDuration.Hours()), int(s.TripDuration.Minutes())%60))
	}

	if s.AlternateDuration != nil {
		fs.AlternateDuration = StringPointer(fmt.Sprintf("%02d%02d", int(s.AlternateDuration.Hours()), int(s.AlternateDuration.Minutes())%60))
	}

	return fs
}
