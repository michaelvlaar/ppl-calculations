package calculations

import (
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/volume"
	"time"
)

const (
	TaxiFuelLiters                    = 2.0
	AquilaFuelLiterConsumptionPerHour = 17.0
	ContingencyFuelTripPercentage     = 0.1
	AquilaReserveFuelLiters           = AquilaFuelLiterConsumptionPerHour * 0.75
)

type FuelPlanning struct {
	Taxi        fuel.Fuel
	Trip        fuel.Fuel
	Alternate   fuel.Fuel
	Contingency fuel.Fuel
	Reserve     fuel.Fuel
	Extra       fuel.Fuel
	Total       fuel.Fuel

	VolumeType volume.Type
	Sufficient bool
}

func NewFuelPlanning(tripDuration time.Duration, alternateDuration time.Duration, f fuel.Fuel, volumeType volume.Type) (*FuelPlanning, error) {
	fp := &FuelPlanning{
		Taxi:        fuel.MustNew(volume.MustNew(TaxiFuelLiters, volume.TypeLiter), f.Type),
		Trip:        fuel.MustNew(volume.MustNew(AquilaFuelLiterConsumptionPerHour*tripDuration.Hours(), volume.TypeLiter), f.Type),
		Alternate:   fuel.MustNew(volume.MustNew(AquilaFuelLiterConsumptionPerHour*alternateDuration.Hours(), volume.TypeLiter), f.Type),
		Contingency: fuel.MustNew(volume.MustNew(AquilaFuelLiterConsumptionPerHour*tripDuration.Hours()*ContingencyFuelTripPercentage, volume.TypeLiter), f.Type),
		Reserve:     fuel.MustNew(volume.MustNew(AquilaReserveFuelLiters, volume.TypeLiter), f.Type),

		VolumeType: volumeType,
	}

	fp.Extra = fuel.Subtract(f, fp.Taxi, fp.Trip, fp.Alternate, fp.Contingency, fp.Reserve)
	fp.Total = fuel.Add(fp.Taxi, fp.Trip, fp.Alternate, fp.Contingency, fp.Reserve, fp.Extra)

	if fp.Extra.Volume.Amount > 0.0 {
		fp.Sufficient = true
	} else {
		fp.Sufficient = false
	}

	return fp, nil
}
