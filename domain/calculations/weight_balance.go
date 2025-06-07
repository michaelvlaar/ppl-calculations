package calculations

import (
	"fmt"
	"github.com/michaelvlaar/ppl-calculations/domain/callsign"
	"github.com/michaelvlaar/ppl-calculations/domain/fuel"
	"github.com/michaelvlaar/ppl-calculations/domain/seat"
	"github.com/michaelvlaar/ppl-calculations/domain/volume"
	"github.com/michaelvlaar/ppl-calculations/domain/weight_balance"
)

const (
	EmptyMassName      = "Empty Mass"
	AquilaEmptyMassArm = 0.4294

	AquilaForwardCgLimit  = 0.427
	AquilaBackwardCgLimit = 0.523
	AquilaMTOW            = 750.0
	AquilaMinWeight       = 558.0
	AquilaMaxFuel         = 109.6

	PilotMassName     = "Pilot"
	PassengerMassName = "Passenger"
	BagageMassName    = "Bagage"
	FuelMassName      = "Fuel"

	AquilaSeatFrontMassArm  = 0.4545
	AquilaSeatMiddleMassArm = 0.5227
	AquilaSeatBackMassArm   = 0.5909

	FuelMassMogasPerLiter = 0.74
	FuelMassAvgasPerLiter = 0.72

	BagageMassArm = 1.3
	FuelMassArm   = 0.3250
)

var (
	AquilaPHDHAEmptyMassMoment = weight_balance.NewMassMoment(EmptyMassName, AquilaEmptyMassArm, weight_balance.NewMass(517.0))
	AquilaPHDHBEmptyMassMoment = weight_balance.NewMassMoment(EmptyMassName, AquilaEmptyMassArm, weight_balance.NewMass(529.5))
)

type WeightBalance struct {
	Moments      []*weight_balance.MassMoment
	Total        *weight_balance.MassMoment
	WithinLimits bool
}

func NewWeightAndBalanceMaxFuel(callSign callsign.CallSign, pilot weight_balance.Mass, pilotSeat seat.Position, passenger *weight_balance.Mass, passengerSeat *seat.Position, bagage weight_balance.Mass, fuelType fuel.Type) (*WeightBalance, fuel.Fuel, error) {
	wb := &WeightBalance{}

	switch callSign.String() {
	case "PHDHA":
		wb.Moments = append(wb.Moments, AquilaPHDHAEmptyMassMoment)
	case "PHDHB":
		wb.Moments = append(wb.Moments, AquilaPHDHBEmptyMassMoment)
	default:
		panic("unknown aircraft")
	}

	var pilotSeatMassArm float64
	switch pilotSeat {
	case seat.PositionFront:
		pilotSeatMassArm = AquilaSeatFrontMassArm
	case seat.PositionMiddle:
		pilotSeatMassArm = AquilaSeatMiddleMassArm
	case seat.PositionBack:
		pilotSeatMassArm = AquilaSeatBackMassArm
	default:
		panic("invalid seat position")
	}

	wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(PilotMassName, pilotSeatMassArm, pilot))

	if passenger == nil || passengerSeat == nil {
		// Add empty mass
		wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(PassengerMassName, pilotSeatMassArm, weight_balance.Mass(0)))
	} else {
		var passengerSeatMassArm float64
		switch pilotSeat {
		case seat.PositionFront:
			passengerSeatMassArm = AquilaSeatFrontMassArm
		case seat.PositionMiddle:
			passengerSeatMassArm = AquilaSeatMiddleMassArm
		default:
			passengerSeatMassArm = AquilaSeatBackMassArm
		}

		wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(PassengerMassName, passengerSeatMassArm, *passenger))
	}

	wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(BagageMassName, BagageMassArm, bagage))

	totalMass := 0.0
	totalKGM := 0.0

	for _, i := range wb.Moments {
		totalMass += i.Mass.Kilo()
		totalKGM += i.KGM()
	}

	calculatedFuelMass := (AquilaForwardCgLimit*totalMass - totalKGM) / (FuelMassArm - AquilaForwardCgLimit)

	var fuelMassPerLiter float64
	switch fuelType {
	case fuel.TypeMogas:
		fuelMassPerLiter = FuelMassMogasPerLiter
	case fuel.TypeAvGas:
		fuelMassPerLiter = FuelMassAvgasPerLiter
	}

	if calculatedFuelMass/fuelMassPerLiter >= AquilaMaxFuel {
		calculatedFuelMass = AquilaMaxFuel * fuelMassPerLiter
	}

	if totalMass+calculatedFuelMass > AquilaMTOW {
		calculatedFuelMass = AquilaMTOW - totalMass
	}

	var f fuel.Fuel
	switch fuelType {
	case fuel.TypeMogas:
		liters := calculatedFuelMass / FuelMassMogasPerLiter
		f = fuel.MustNew(volume.MustNew(liters, volume.TypeLiter), fuel.TypeMogas)
	case fuel.TypeAvGas:
		liters := calculatedFuelMass / FuelMassAvgasPerLiter
		f = fuel.MustNew(volume.MustNew(liters, volume.TypeLiter), fuel.TypeAvGas)
	}

	var volumeDescription string
	var fuelMass weight_balance.Mass
	switch f.Volume.Type {
	case volume.TypeLiter:
		fuelMass = weight_balance.NewMass(f.Volume.Amount * fuelMassPerLiter)
		volumeDescription = volume.DescriptionLiter
	case volume.TypeGallon:
		fuelMass = weight_balance.NewMass(f.Volume.Amount * volume.LitersInGallon * fuelMassPerLiter)
		volumeDescription = volume.DescriptionGallon
	}

	wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(fmt.Sprintf("%s (%.2fkg/%s)", FuelMassName, fuelMassPerLiter, volumeDescription), FuelMassArm, fuelMass))

	totalMass = 0.0
	totalKGM = 0.0
	for _, i := range wb.Moments {
		totalMass += i.Mass.Kilo()
		totalKGM += i.KGM()
	}

	wb.Total = weight_balance.NewMassMoment("Total", totalKGM/totalMass, weight_balance.NewMass(totalMass))

	wb.WithinLimits = totalMass <= AquilaMTOW && totalKGM/totalMass >= AquilaForwardCgLimit && totalKGM/totalMass <= AquilaBackwardCgLimit

	return wb, f, nil
}

func NewWeightAndBalance(callSign callsign.CallSign, pilot weight_balance.Mass, pilotSeat seat.Position, passenger *weight_balance.Mass, passengerSeat *seat.Position, bagage *weight_balance.Mass, f fuel.Fuel) (*WeightBalance, error) {
	wb := &WeightBalance{}

	switch callSign.String() {
	case "PHDHA":
		wb.Moments = append(wb.Moments, AquilaPHDHAEmptyMassMoment)
	case "PHDHB":
		wb.Moments = append(wb.Moments, AquilaPHDHBEmptyMassMoment)
	default:
		panic("unknown aircraft")
	}

	var pilotSeatMassArm float64
	switch pilotSeat {
	case seat.PositionFront:
		pilotSeatMassArm = AquilaSeatFrontMassArm
	case seat.PositionMiddle:
		pilotSeatMassArm = AquilaSeatMiddleMassArm
	case seat.PositionBack:
		pilotSeatMassArm = AquilaSeatBackMassArm
	default:
		panic("invalid seat position")
	}

	wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(PilotMassName, pilotSeatMassArm, pilot))

	if passenger == nil || passengerSeat == nil {
		wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(PassengerMassName, pilotSeatMassArm, weight_balance.Mass(0)))
	} else {
		var passengerSeatMassArm float64
		switch pilotSeat {
		case seat.PositionFront:
			passengerSeatMassArm = AquilaSeatFrontMassArm
		case seat.PositionMiddle:
			passengerSeatMassArm = AquilaSeatMiddleMassArm
		default:
			passengerSeatMassArm = AquilaSeatBackMassArm
		}

		wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(PassengerMassName, passengerSeatMassArm, *passenger))
	}

	if bagage != nil {
		wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(BagageMassName, BagageMassArm, *bagage))
	} else {
		wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(BagageMassName, BagageMassArm, weight_balance.NewMass(0.0)))
	}

	var fuelMassPerLiter float64
	switch f.Type {
	case fuel.TypeMogas:
		fuelMassPerLiter = FuelMassMogasPerLiter
	case fuel.TypeAvGas:
		fuelMassPerLiter = FuelMassAvgasPerLiter
	default:
		panic("invalid fuel type")
	}

	var fuelMass weight_balance.Mass
	var volumeDescription string
	switch f.Volume.Type {
	case volume.TypeLiter:
		fuelMass = weight_balance.NewMass(f.Volume.Amount * fuelMassPerLiter)
		volumeDescription = volume.DescriptionLiter
	case volume.TypeGallon:
		fuelMass = weight_balance.NewMass(f.Volume.Amount * volume.LitersInGallon * fuelMassPerLiter)
		volumeDescription = volume.DescriptionGallon
	}

	wb.Moments = append(wb.Moments, weight_balance.NewMassMoment(fmt.Sprintf("%s (%.2fkg/%s)", FuelMassName, fuelMassPerLiter, volumeDescription), FuelMassArm, fuelMass))

	totalMass := 0.0
	totalKGM := 0.0

	for _, i := range wb.Moments {
		totalMass += i.Mass.Kilo()
		totalKGM += i.KGM()
	}

	wb.Total = weight_balance.NewMassMoment("Total", totalKGM/totalMass, weight_balance.NewMass(totalMass))

	wb.WithinLimits = totalMass <= AquilaMTOW && totalKGM/totalMass >= AquilaForwardCgLimit && totalKGM/totalMass <= AquilaBackwardCgLimit

	return wb, nil
}
