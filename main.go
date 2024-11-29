package main

import (
	"bytes"
	"context"
	"embed"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"ppl-calculations/adapters"
	"ppl-calculations/app"
	"ppl-calculations/app/commands"
	"ppl-calculations/app/queries"
	"ppl-calculations/ports"
	"sync"
)

//go:embed assets/*
var assets embed.FS

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("application started")
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ldrFile, err := assets.Open("assets/charts/ldr.svg")
	if err != nil {
		panic(err)
	}

	var ldrBuf bytes.Buffer
	_, err = io.Copy(&ldrBuf, ldrFile)
	if err != nil {
		panic(err)
	}

	err = ldrFile.Close()
	if err != nil {
		panic(err)
	}

	tdrFile, err := assets.Open("assets/charts/tdr.svg")
	if err != nil {
		panic(err)
	}

	var tdrBuf bytes.Buffer
	_, err = io.Copy(&tdrBuf, tdrFile)
	if err != nil {
		panic(err)
	}

	err = tdrFile.Close()
	if err != nil {
		panic(err)
	}

	calculationsService := adapters.NewCalculationsService(tdrBuf, ldrBuf)

	a := app.Application{
		Commands: app.Commands{
			UpdateLoadSheet: commands.NewUpdateLoadSheetHandler(),
			UpdateFuelSheet: commands.NewUpdateFuelSheetHandler(),
		},
		Queries: app.Queries{
			WBChart:    queries.NewWBChartHandler(),
			LoadSheet:  queries.NewLoadSheetHandler(),
			FuelSheet:  queries.NewFuelSheetHandler(),
			StatsSheet: queries.NewStatsSheetHandler(calculationsService),
			LdrChart:   queries.NewLdrChartHandler(calculationsService),
			TodChart:   queries.NewTodChartHandler(calculationsService),
			PdfExport:  queries.NewPdfExportHandler(),
		},
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ports.NewHTTPListener(ctx, &wg, a, assets)
	}()

	<-stop
	logrus.Info("graceful shutdown received")

	cancel()
	wg.Wait()

	logrus.Info("shutting down")
}
