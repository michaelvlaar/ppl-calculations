package main

import (
	"context"
	"embed"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"os/signal"
	"ppl-calculations/adapters"
	"ppl-calculations/app"
	"ppl-calculations/app/commands"
	"ppl-calculations/app/queries"
	"ppl-calculations/ports/http"
	"sync"
)

//go:embed assets/*
var assets embed.FS

var version = "dev"

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("application started")
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	chartsFS, err := fs.Sub(assets, "assets/charts")
	if err != nil {
		logrus.WithError(err).Fatal("could not load chart folder")
	}

	calculationsService := adapters.MustNewCalculationsService(chartsFS, adapters.MustNewImageService())

	exportFS, err := fs.Sub(assets, "assets/export")
	if err != nil {
		logrus.WithError(err).Fatal("could not load export folder")
	}

	a := app.Application{
		Commands: app.Commands{
			UpdateLoadSheet:   commands.NewUpdateLoadSheetHandler(),
			UpdateFuelSheet:   commands.NewUpdateFuelSheetHandler(),
			UpdateExportSheet: commands.NewUpdateExportSheetHandler(),
			DeleteExportSheet: commands.NewDeleteExportSheetHandler(),
			ClearSheet:        commands.NewClearSheetHandler(),
		},
		Queries: app.Queries{
			WBChart:     queries.NewWBChartHandler(calculationsService),
			LoadSheet:   queries.NewLoadSheetHandler(),
			FuelSheet:   queries.NewFuelSheetHandler(),
			StatsSheet:  queries.NewStatsSheetHandler(calculationsService),
			ExportSheet: queries.NewExportSheetHandler(),
			Export:      queries.NewExportHandler(),
			Exports:     queries.NewExportsHandler(),
			LdrChart:    queries.NewLdrChartHandler(calculationsService),
			TodChart:    queries.NewTodChartHandler(calculationsService),
			PdfExport:   queries.NewPdfExportHandler(exportFS, calculationsService),
		},
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.NewHTTPListener(ctx, &wg, a, assets, version)
	}()

	<-stop
	logrus.Info("graceful shutdown received")

	cancel()
	wg.Wait()

	logrus.Info("shutting down")
}
