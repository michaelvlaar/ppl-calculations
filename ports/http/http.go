package http

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"io/fs"
	"net/http"
	"os"
	"ppl-calculations/app"
	"ppl-calculations/ports/http/middleware"
	"ppl-calculations/ports/http/routes"
	"ppl-calculations/ports/templates"
	"sync"
	"time"
)

func NewHTTPListener(ctx context.Context, wg *sync.WaitGroup, app app.Application, assets fs.FS, version string) {
	mux := http.NewServeMux()

	routes.RegisterStaticRoutes(mux, assets)
	routes.RegisterHealthRoutes(mux)
	routes.RegisterOverviewRoutes(mux, app)
	routes.RegisterCalculationRoutes(mux, app)
	routes.RegisterChartRoutes(mux, app)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: middleware.Chain(mux, middleware.SecurityHeaders, templates.HttpMiddleware(version), middleware.CSRF),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.WithField("addr", server.Addr).Info("starting HTTP server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Error("error starting HTTP server")
		}
		logrus.Info("HTTP server closed")
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logrus.WithError(err).Error("error shutting down HTTP server")
		}
	}
}
