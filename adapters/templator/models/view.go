package models

import (
	"fmt"
	"ppl-calculations/app/queries"
	"ppl-calculations/domain/export"
	"ppl-calculations/domain/wind"
	"strings"
)

type View struct {
	Base

	Name string
	Date string

	FuelTaxi        string
	FuelTrip        string
	FuelAlternate   string
	FuelContingency string
	FuelReserve     string
	FuelTotal       string
	FuelExtra       string
	FuelExtraAbs    string
	FuelSufficient  bool

	ChartUrl string
	LdrUrl   string
	TdrUrl   string

	TakeOffRunRequired        string
	TakeOffDistanceRequired   string
	LandingDistanceRequired   string
	LandingGroundRollRequired string

	WeightAndBalanceTakeOff WeightAndBalanceState
	WeightAndBalanceLanding WeightAndBalanceState
}

func ViewFromExport(csrf string, statsSheet queries.StatsSheetResponse, e export.Export) interface{} {
	template := View{
		Base: Base{
			Step: string(StepView),
			CSRF: csrf,
		},
		Name: e.Name.String(),
		Date: e.CreatedAt.Format("15:04:05 02-01-2006"),
	}

	template.FuelSufficient = statsSheet.Calculations.FuelPlanning.Sufficient
	template.FuelTaxi = parseNumber(statsSheet.Calculations.FuelPlanning.Taxi.Volume.String(statsSheet.Calculations.FuelPlanning.VolumeType))
	template.FuelTrip = parseNumber(statsSheet.Calculations.FuelPlanning.Trip.Volume.String(statsSheet.Calculations.FuelPlanning.VolumeType))
	template.FuelAlternate = parseNumber(statsSheet.Calculations.FuelPlanning.Alternate.Volume.String(statsSheet.Calculations.FuelPlanning.VolumeType))
	template.FuelContingency = parseNumber(statsSheet.Calculations.FuelPlanning.Contingency.Volume.String(statsSheet.Calculations.FuelPlanning.VolumeType))
	template.FuelReserve = parseNumber(statsSheet.Calculations.FuelPlanning.Reserve.Volume.String(statsSheet.Calculations.FuelPlanning.VolumeType))
	template.FuelTotal = parseNumber(statsSheet.Calculations.FuelPlanning.Total.Volume.String(statsSheet.Calculations.FuelPlanning.VolumeType))
	template.FuelExtra = parseNumber(statsSheet.Calculations.FuelPlanning.Extra.Volume.String(statsSheet.Calculations.FuelPlanning.VolumeType))
	template.FuelExtraAbs = parseNumber(statsSheet.Calculations.FuelPlanning.Extra.Volume.Abs().String(statsSheet.Calculations.FuelPlanning.VolumeType))

	wbState := WeightAndBalanceState{}
	for _, i := range statsSheet.Calculations.TakeOffWeightAndBalance.Moments {
		m := parseNumber(fmt.Sprintf("%.2f", i.Mass.Kilo()))
		if strings.HasPrefix(i.Name, "Fuel") {
			m = fmt.Sprintf("(%s) %s", template.FuelTotal, m)
		}

		wbState.Items = append(wbState.Items, WeightAndBalanceItem{
			Name:       parseNumber(i.Name),
			LeverArm:   parseNumber(fmt.Sprintf("%.4f", i.Arm)),
			Mass:       m,
			MassMoment: parseNumber(fmt.Sprintf("%.2f", i.KGM())),
		})
	}

	wbState.Total = WeightAndBalanceItem{
		Name:       parseNumber(statsSheet.Calculations.TakeOffWeightAndBalance.Total.Name),
		LeverArm:   parseNumber(fmt.Sprintf("%.4f", statsSheet.Calculations.TakeOffWeightAndBalance.Total.Arm)),
		Mass:       parseNumber(fmt.Sprintf("%.2f", statsSheet.Calculations.TakeOffWeightAndBalance.Total.Mass.Kilo())),
		MassMoment: parseNumber(fmt.Sprintf("%.2f", statsSheet.Calculations.TakeOffWeightAndBalance.Total.KGM())),
	}

	wbState.WithinLimits = statsSheet.Calculations.TakeOffWeightAndBalance.WithinLimits

	wbLandingState := WeightAndBalanceState{}

	for _, i := range statsSheet.Calculations.LandingWeightAndBalance.Moments {
		m := parseNumber(fmt.Sprintf("%.2f", i.Mass.Kilo()))
		if strings.HasPrefix(i.Name, "Fuel") {
			m = fmt.Sprintf("(%s) %s", parseNumber(statsSheet.Calculations.FuelPlanning.Total.Volume.Subtract(statsSheet.Calculations.FuelPlanning.Trip.Volume).String(statsSheet.Calculations.FuelPlanning.VolumeType)), m)
		}

		wbLandingState.Items = append(wbLandingState.Items, WeightAndBalanceItem{
			Name:       parseNumber(i.Name),
			LeverArm:   parseNumber(fmt.Sprintf("%.4f", i.Arm)),
			Mass:       m,
			MassMoment: parseNumber(fmt.Sprintf("%.2f", i.KGM())),
		})
	}

	wbLandingState.Total = WeightAndBalanceItem{
		Name:       parseNumber(statsSheet.Calculations.LandingWeightAndBalance.Total.Name),
		LeverArm:   parseNumber(fmt.Sprintf("%.4f", statsSheet.Calculations.LandingWeightAndBalance.Total.Arm)),
		Mass:       parseNumber(fmt.Sprintf("%.2f", statsSheet.Calculations.LandingWeightAndBalance.Total.Mass.Kilo())),
		MassMoment: parseNumber(fmt.Sprintf("%.2f", statsSheet.Calculations.LandingWeightAndBalance.Total.KGM())),
	}

	wbLandingState.WithinLimits = statsSheet.Calculations.TakeOffWeightAndBalance.WithinLimits

	template.ChartUrl = fmt.Sprintf("/aquila-wb?callsign=%s&takeoff-mass=%.2f&takeoff-mass-moment=%.2f&landing-mass=%.2f&landing-mass-moment=%.2f&limits=%t", statsSheet.Calculations.CallSign.String(), statsSheet.Calculations.TakeOffWeightAndBalance.Total.Mass, statsSheet.Calculations.TakeOffWeightAndBalance.Total.KGM(), statsSheet.Calculations.LandingWeightAndBalance.Total.Mass, statsSheet.Calculations.LandingWeightAndBalance.Total.KGM(), statsSheet.Calculations.TakeOffWeightAndBalance.WithinLimits)

	wd := "headwind"
	if statsSheet.Calculations.Wind.Direction == wind.DirectionTailwind {
		wd = "tailwind"
	}

	template.LdrUrl = fmt.Sprintf("/aquila-ldr?pressure_altitude=%0.2f&mtow=%.2f&oat=%.2f&wind=%.2f&wind_direction=%s", statsSheet.Calculations.PressureAltitude.Float64(), statsSheet.Calculations.LandingWeightAndBalance.Total.Mass, statsSheet.Calculations.OAT.Float64(), statsSheet.Calculations.Wind.Speed.Float64(), wd)
	template.TdrUrl = fmt.Sprintf("/aquila-tdr?pressure_altitude=%0.2f&mtow=%.2f&oat=%.2f&wind=%.2f&wind_direction=%s", statsSheet.Calculations.PressureAltitude.Float64(), statsSheet.Calculations.TakeOffWeightAndBalance.Total.Mass, statsSheet.Calculations.OAT.Float64(), statsSheet.Calculations.Wind.Speed.Float64(), wd)

	template.TakeOffDistanceRequired = fmt.Sprintf("%.0f", statsSheet.Calculations.Performance.TakeOffDistanceRequired)
	template.TakeOffRunRequired = fmt.Sprintf("%.0f", statsSheet.Calculations.Performance.TakeOffRunRequired)
	template.LandingDistanceRequired = fmt.Sprintf("%.0f", statsSheet.Calculations.Performance.LandingDistanceRequired)
	template.LandingGroundRollRequired = fmt.Sprintf("%.0f", statsSheet.Calculations.Performance.LandingGroundRollRequired)

	template.WeightAndBalanceTakeOff = wbState
	template.WeightAndBalanceLanding = wbLandingState

	return template
}
