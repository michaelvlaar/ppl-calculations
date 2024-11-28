package app

import (
	"ppl-calculations/app/commands"
	"ppl-calculations/app/queries"
)

type Application struct {
	Queries  Queries
	Commands Commands
}

type Commands struct {
	UpdateLoadSheet commands.UpdateLoadSheetHandler
	UpdateFuelSheet commands.UpdateFuelSheetHandler
}

type Queries struct {
	LoadSheet  queries.LoadSheetHandler
	FuelSheet  queries.FuelSheetHandler
	StatsSheet queries.StatsSheetHandler
	WBChart    queries.WBChartHandler
}
