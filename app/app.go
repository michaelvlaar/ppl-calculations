package app

import (
	"github.com/michaelvlaar/ppl-calculations/app/commands"
	"github.com/michaelvlaar/ppl-calculations/app/queries"
)

type Application struct {
	Queries  Queries
	Commands Commands
}

type Commands struct {
	UpdateLoadSheet   commands.UpdateLoadSheetHandler
	UpdateFuelSheet   commands.UpdateFuelSheetHandler
	UpdateExportSheet commands.UpdateExportSheetHandler
	DeleteExportSheet commands.DeleteExportSheetHandler
	ClearSheet        commands.ClearSheetHandler
}

type Queries struct {
	LoadSheet   queries.LoadSheetHandler
	FuelSheet   queries.FuelSheetHandler
	StatsSheet  queries.StatsSheetHandler
	ExportSheet queries.ExportSheetHandler

	Exports queries.ExportsHandler

	WBChart  queries.WBChartHandler
	LdrChart queries.LdrChartHandler
	TodChart queries.TodChartHandler

	PdfExport queries.PdfExportHandler
}
