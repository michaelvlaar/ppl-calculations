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

	Total fuel.Fuel

	Sufficient bool
}

func NewMaxFuelPlanning(fuelType fuel.Type, tripDuration time.Duration, alternateDuration time.Duration) (*FuelPlanning, error) {
	fp := &FuelPlanning{
		Taxi:        fuel.MustNew(volume.MustNew(TaxiFuelLiters, volume.TypeLiter), fuelType),
		Trip:        fuel.MustNew(volume.MustNew(AquilaFuelLiterConsumptionPerHour*tripDuration.Hours(), volume.TypeLiter), fuelType),
		Alternate:   fuel.MustNew(volume.MustNew(AquilaFuelLiterConsumptionPerHour*alternateDuration.Hours(), volume.TypeLiter), fuelType),
		Contingency: fuel.MustNew(volume.MustNew(AquilaFuelLiterConsumptionPerHour*tripDuration.Hours()*ContingencyFuelTripPercentage, volume.TypeLiter), fuelType),
		Reserve:     fuel.MustNew(volume.MustNew(AquilaReserveFuelLiters, volume.TypeLiter), fuelType),
		// TODO: calculate max fuel based on weight and balance maximum
		Extra: fuel.MustNew(volume.MustNew(0, volume.TypeLiter), fuelType),
	}

	fp.Total = fuel.Add(fp.Taxi, fp.Trip, fp.Alternate, fp.Contingency, fp.Reserve, fp.Extra)

	fp.Sufficient = false

	return fp, nil
}
