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
	"ppl-calculations/ports/http/middleware"
	"ppl-calculations/ports/templates"
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
	stateServiceProvider := adapters.NewCookieStateServiceProvider()

	exportFS, err := fs.Sub(assets, "assets/export")
	if err != nil {
		logrus.WithError(err).Fatal("could not load export folder")
	}

	a := app.Application{
		Commands: app.Commands{
			UpdateLoadSheet:   commands.NewUpdateLoadSheetHandler(stateServiceProvider),
			UpdateFuelSheet:   commands.NewUpdateFuelSheetHandler(stateServiceProvider),
			UpdateExportSheet: commands.NewUpdateExportSheetHandler(stateServiceProvider),
			DeleteExportSheet: commands.NewDeleteExportSheetHandler(stateServiceProvider),
			ClearSheet:        commands.NewClearSheetHandler(stateServiceProvider),
		},
		Queries: app.Queries{
			WBChart:     queries.NewWBChartHandler(calculationsService),
			LoadSheet:   queries.NewLoadSheetHandler(stateServiceProvider),
			FuelSheet:   queries.NewFuelSheetHandler(stateServiceProvider),
			StatsSheet:  queries.NewStatsSheetHandler(stateServiceProvider, calculationsService),
			ExportSheet: queries.NewExportSheetHandler(stateServiceProvider),
			Exports:     queries.NewExportsHandler(stateServiceProvider),
			LdrChart:    queries.NewLdrChartHandler(calculationsService),
			TodChart:    queries.NewTodChartHandler(calculationsService),
			PdfExport:   queries.NewPdfExportHandler(exportFS, calculationsService),
		},
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.NewHTTPListener(ctx, &wg, a, assets, middleware.SecurityHeaders, templates.HttpMiddleware(version), adapters.HttpMiddleware(stateServiceProvider))
	}()

	<-stop
	logrus.Info("graceful shutdown received")

	cancel()
	wg.Wait()

	logrus.Info("shutting down")
}
