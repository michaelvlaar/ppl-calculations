package main

import (
	"context"
	"embed"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
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

	a := app.Application{
		Commands: app.Commands{
			UpdateLoadSheet: commands.NewUpdateLoadSheetHandler(),
			UpdateFuelSheet: commands.NewUpdateFuelSheetHandler(),
		},
		Queries: app.Queries{
			WBChart:   queries.NewWBChartHandler(),
			LoadSheet: queries.NewLoadSheetHandler(),
			FuelSheet: queries.NewFuelSheetHandler(),
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
