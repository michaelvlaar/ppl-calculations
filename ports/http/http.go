package http

import (
	"context"
	"errors"
	"github.com/michaelvlaar/ppl-calculations/app"
	"github.com/michaelvlaar/ppl-calculations/ports/http/middleware"
	"github.com/michaelvlaar/ppl-calculations/ports/http/routes"
	"github.com/nanmu42/gzip"
	"github.com/sirupsen/logrus"
	"io/fs"
	"net/http"
	"os"
	"sync"
	"time"
)

func NewHTTPListener(ctx context.Context, wg *sync.WaitGroup, app app.Application, assets fs.FS, middlewares ...func(http.Handler) http.Handler) {
	mux := http.NewServeMux()

	routes.RegisterStaticRoutes(mux, assets)
	routes.RegisterHealthRoutes(mux)
	routes.RegisterOverviewRoutes(mux, app)
	routes.RegisterCalculationRoutes(mux, app)
	routes.RegisterChartRoutes(mux, app)

	middlewares = append(middlewares, middleware.CSRF)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: middleware.Chain(gzip.DefaultHandler().WrapHandler(mux), middlewares...),
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
