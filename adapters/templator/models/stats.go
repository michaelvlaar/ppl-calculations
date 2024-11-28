package models

import (
	"fmt"
	"ppl-calculations/app/queries"
	"strings"
)

type WeightAndBalanceItem struct {
	Name       string
	LeverArm   string
	Mass       string
	MassMoment string
}

type WeightAndBalanceState struct {
	Items        []WeightAndBalanceItem
	Total        WeightAndBalanceItem
	WithinLimits bool
}

type Stats struct {
	Base

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

	WeightAndBalanceTakeOff WeightAndBalanceState
	WeightAndBalanceLanding WeightAndBalanceState
}

func parseNumber(number string) string {
	return strings.ReplaceAll(number, ".", ",")
}

func StatsFromStatsSheet(statsSheet queries.StatsSheetResponse) interface{} {
	template := Stats{
		Base: Base{
			Step: string(StepStats),
		},
	}

	template.FuelSufficient = statsSheet.FuelPlanning.Sufficient
	template.FuelTaxi = parseNumber(statsSheet.FuelPlanning.Taxi.Volume.String(statsSheet.FuelPlanning.VolumeType))
	template.FuelTrip = parseNumber(statsSheet.FuelPlanning.Trip.Volume.String(statsSheet.FuelPlanning.VolumeType))
	template.FuelAlternate = parseNumber(statsSheet.FuelPlanning.Alternate.Volume.String(statsSheet.FuelPlanning.VolumeType))
	template.FuelContingency = parseNumber(statsSheet.FuelPlanning.Contingency.Volume.String(statsSheet.FuelPlanning.VolumeType))
	template.FuelReserve = parseNumber(statsSheet.FuelPlanning.Reserve.Volume.String(statsSheet.FuelPlanning.VolumeType))
	template.FuelTotal = parseNumber(statsSheet.FuelPlanning.Total.Volume.String(statsSheet.FuelPlanning.VolumeType))
	template.FuelExtra = parseNumber(statsSheet.FuelPlanning.Extra.Volume.String(statsSheet.FuelPlanning.VolumeType))
	template.FuelExtraAbs = parseNumber(statsSheet.FuelPlanning.Extra.Volume.String(statsSheet.FuelPlanning.VolumeType))

	wbState := WeightAndBalanceState{}
	for _, i := range statsSheet.TakeOffWeightAndBalance.Moments {
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
		Name:       parseNumber(statsSheet.TakeOffWeightAndBalance.Total.Name),
		LeverArm:   parseNumber(fmt.Sprintf("%.4f", statsSheet.TakeOffWeightAndBalance.Total.Arm)),
		Mass:       parseNumber(fmt.Sprintf("%.2f", statsSheet.TakeOffWeightAndBalance.Total.Mass.Kilo())),
		MassMoment: parseNumber(fmt.Sprintf("%.2f", statsSheet.TakeOffWeightAndBalance.Total.KGM())),
	}

	wbState.WithinLimits = statsSheet.TakeOffWeightAndBalance.WithinLimits

	wbLandingState := WeightAndBalanceState{}

	for _, i := range statsSheet.LandingWeightAndBalance.Moments {
		m := parseNumber(fmt.Sprintf("%.2f", i.Mass.Kilo()))
		if strings.HasPrefix(i.Name, "Fuel") {
			m = fmt.Sprintf("(%s) %s", parseNumber(statsSheet.FuelPlanning.Total.Volume.Subtract(statsSheet.FuelPlanning.Trip.Volume).String(statsSheet.FuelPlanning.VolumeType)), m)
		}

		wbLandingState.Items = append(wbLandingState.Items, WeightAndBalanceItem{
			Name:       parseNumber(i.Name),
			LeverArm:   parseNumber(fmt.Sprintf("%.4f", i.Arm)),
			Mass:       m,
			MassMoment: parseNumber(fmt.Sprintf("%.2f", i.KGM())),
		})
	}

	wbLandingState.Total = WeightAndBalanceItem{
		Name:       parseNumber(statsSheet.LandingWeightAndBalance.Total.Name),
		LeverArm:   parseNumber(fmt.Sprintf("%.4f", statsSheet.LandingWeightAndBalance.Total.Arm)),
		Mass:       parseNumber(fmt.Sprintf("%.2f", statsSheet.LandingWeightAndBalance.Total.Mass.Kilo())),
		MassMoment: parseNumber(fmt.Sprintf("%.2f", statsSheet.LandingWeightAndBalance.Total.KGM())),
	}

	wbState.WithinLimits = statsSheet.TakeOffWeightAndBalance.WithinLimits

	template.ChartUrl = fmt.Sprintf("/aquila-wb?callsign=%s&takeoff-mass=%.2f&takeoff-mass-moment=%.2f&landing-mass=%.2f&landing-mass-moment=%.2f&limits=%t", statsSheet.CallSign.String(), statsSheet.TakeOffWeightAndBalance.Total.Mass, statsSheet.TakeOffWeightAndBalance.Total.KGM(), statsSheet.LandingWeightAndBalance.Total.Mass, statsSheet.LandingWeightAndBalance.Total.KGM(), statsSheet.TakeOffWeightAndBalance.WithinLimits)

	template.WeightAndBalanceTakeOff = wbState
	template.WeightAndBalanceLanding = wbLandingState

	return template
}
