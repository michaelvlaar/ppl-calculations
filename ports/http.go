package ports

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"ppl-calculations/adapters"
	"ppl-calculations/adapters/templator/models"
	"ppl-calculations/adapters/templator/parsing"
	"ppl-calculations/app"
	"ppl-calculations/app/queries"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func NewHTTPListener(ctx context.Context, wg *sync.WaitGroup, app app.Application, assets fs.FS) {
	mux := http.NewServeMux()

	templatesFS, err := fs.Sub(assets, "assets/templates")
	if err != nil {
		logrus.WithError(err).Fatal("cannot open asset templates")
	}

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"derefString": derefString,
		"mod": func(i, j int) bool {
			return i%j != 0
		},
	}).ParseFS(templatesFS, "*.html")
	if err != nil {
		logrus.WithError(err).Fatal("cannot parse asset templates")
	}

	cssFs, err := fs.Sub(assets, "assets/css")
	if err != nil {
		log.Fatalf("Fout bij het parsen van css: %v", err)
	}
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.FS(cssFs))))

	jsFs, err := fs.Sub(assets, "assets/js")
	if err != nil {
		log.Fatalf("Fout bij het parsen van css: %v", err)
	}
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.FS(jsFs))))

	fontFs, err := fs.Sub(assets, "assets/fonts")
	if err != nil {
		log.Fatalf("Fout bij het parsen van css: %v", err)
	}
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.FS(fontFs))))

	mux.HandleFunc("/aquila-wb", func(w http.ResponseWriter, r *http.Request) {
		urlTakeOffMass := r.URL.Query().Get("takeoff-mass")
		urlTakeOffMassMoment := r.URL.Query().Get("takeoff-mass-moment")
		urlLandingMass := r.URL.Query().Get("landing-mass")
		urlLandingMassMoment := r.URL.Query().Get("landing-mass-moment")
		urlLimits := r.URL.Query().Get("limits")
		urlCallSign := r.URL.Query().Get("callsign")

		if urlTakeOffMass == "" || urlTakeOffMassMoment == "" || urlLandingMass == "" || urlLandingMassMoment == "" || urlLimits == "" || urlCallSign == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		takeOffMass, err := strconv.ParseFloat(urlTakeOffMass, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		takeOffMassMoment, err := strconv.ParseFloat(urlTakeOffMassMoment, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		landingMass, err := strconv.ParseFloat(urlLandingMass, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		landingMassMoment, err := strconv.ParseFloat(urlLandingMassMoment, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		limits := false
		if urlLimits == "true" {
			limits = true
		}

		cs, err := callsign.New(urlCallSign)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chart, err := app.Queries.WBChart.Handle(r.Context(), queries.WBChartRequest{
			CallSign:          cs,
			TakeOffMassMoment: *weight_balance.NewMassMoment("Take-off", takeOffMassMoment/takeOffMass, weight_balance.NewMass(takeOffMass)),
			LandingMassMoment: *weight_balance.NewMassMoment("Landing", landingMassMoment/landingMass, weight_balance.NewMass(landingMass)),
			WithinLimits:      limits,
			ChartType:         calculations.ChartTypeSVG,
		})

		w.Header().Set("Content-Type", "image/svg+xml")

		_, err = io.Copy(w, chart)
		if err != nil {
			logrus.WithError(err).Error("writing chart")
		}
	})

	mux.HandleFunc("/aquila-tdr", func(w http.ResponseWriter, r *http.Request) {
		urlPa := r.URL.Query().Get("pressure_altitude")
		urlOAT := r.URL.Query().Get("oat")
		urlMTOW := r.URL.Query().Get("mtow")
		urlWind := r.URL.Query().Get("wind")
		urlWindDirection := r.URL.Query().Get("wind_direction")

		if urlPa == "" || urlOAT == "" || urlMTOW == "" || urlWind == "" || urlWindDirection == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pa, err := pressure.NewFromString(urlPa)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		oat, err := temperature.NewFromString(urlOAT)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mtow, err := strconv.ParseFloat(urlMTOW, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tow := weight_balance.NewMass(mtow)

		d, err := wind.NewDirectionFromString(urlWindDirection)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s, err := wind.NewSpeedFromString(urlWind)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		wi, err := wind.New(d, s)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chart, err := app.Queries.TodChart.Handle(r.Context(), queries.TodChartRequest{
			OAT:              oat,
			PressureAltitude: pa,
			Tow:              tow,
			Wind:             wi,
		})
		if err != nil {
			logrus.WithError(err).Error("creating chart")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")

		_, err = io.Copy(w, chart.Chart)
		if err != nil {
			logrus.WithError(err).Error("writing chart")
		}
	})

	mux.HandleFunc("/aquila-ldr", func(w http.ResponseWriter, r *http.Request) {
		urlPa := r.URL.Query().Get("pressure_altitude")
		urlOAT := r.URL.Query().Get("oat")
		urlMTOW := r.URL.Query().Get("mtow")
		urlWind := r.URL.Query().Get("wind")
		urlWindDirection := r.URL.Query().Get("wind_direction")

		if urlPa == "" || urlOAT == "" || urlMTOW == "" || urlWind == "" || urlWindDirection == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pa, err := pressure.NewFromString(urlPa)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		oat, err := temperature.NewFromString(urlOAT)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mtow, err := strconv.ParseFloat(urlMTOW, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tow := weight_balance.NewMass(mtow)

		d, err := wind.NewDirectionFromString(urlWindDirection)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s, err := wind.NewSpeedFromString(urlWind)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		wi, err := wind.New(d, s)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		chart, err := app.Queries.LdrChart.Handle(r.Context(), queries.LdrChartRequest{
			OAT:              oat,
			PressureAltitude: pa,
			Tow:              tow,
			Wind:             wi,
		})
		if err != nil {
			logrus.WithError(err).Error("creating chart")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")

		_, err = io.Copy(w, chart.Chart)
		if err != nil {
			logrus.WithError(err).Error("writing chart")
		}
	})

	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		stateService, err := adapters.NewCookieStateService(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			if r.Header.Get("HX-Request") != "true" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = app.Commands.ClearSheet.Handle(r.Context(), stateService); err != nil {
				logrus.WithError(err).Error("clear sheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			loadSheet, err := app.Queries.LoadSheet.Handle(r.Context(), stateService)
			if err != nil {
				logrus.WithError(err).Error("reading loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Push-Url", "/")
			if err = tmpl.ExecuteTemplate(w, "wb_form", models.WeightFromLoadSheet(csrf.Token(r), loadSheet)); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		stateService, err := adapters.NewCookieStateService(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			if r.Header.Get("HX-Request") != "true" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			updateLoadSheetRequest, err := parsing.UpdateLoadSheetRequest(r)
			if err != nil {
				logrus.WithError(err).Error("creating update loadsheet request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = app.Commands.UpdateLoadSheet.Handle(r.Context(), stateService, updateLoadSheetRequest); err != nil {
				logrus.WithError(err).Error("update loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fuelSheet, err := app.Queries.FuelSheet.Handle(r.Context(), stateService)
			if err != nil {
				logrus.WithError(err).Error("reading fuelsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Push-Url", "/fuel")
			if err = tmpl.ExecuteTemplate(w, "fuel_form", models.FuelFromFuelSheet(csrf.Token(r), fuelSheet)); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		case http.MethodGet:
			loadSheet, err := app.Queries.LoadSheet.Handle(r.Context(), stateService)
			if err != nil {
				logrus.WithError(err).Error("reading loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "weight" {
				w.Header().Set("HX-Push-Url", "/")
				if err = tmpl.ExecuteTemplate(w, "wb_form", models.WeightFromLoadSheet(csrf.Token(r), loadSheet)); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = tmpl.ExecuteTemplate(w, "index.html", models.WeightFromLoadSheet(csrf.Token(r), loadSheet)); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/fuel", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		stateService, err := adapters.NewCookieStateService(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			if r.Header.Get("HX-Request") != "true" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			updateFuelSheetRequest, err := parsing.UpdateFuelSheetRequest(r)
			if err != nil {
				logrus.WithError(err).Error("creating update fuelsheet request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = app.Commands.UpdateFuelSheet.Handle(r.Context(), stateService, updateFuelSheetRequest); err != nil {
				logrus.WithError(err).Error("update loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			statsSheet, err := app.Queries.StatsSheet.Handle(r.Context(), stateService)
			if err != nil {
				logrus.WithError(err).Error("update loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Push-Url", "/stats")
			if err = tmpl.ExecuteTemplate(w, "calculations_form", models.StatsFromStatsSheet(csrf.Token(r), statsSheet)); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		case http.MethodGet:
			fuelSheet, err := app.Queries.FuelSheet.Handle(r.Context(), stateService)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "fuel" {
				w.Header().Set("HX-Push-Url", "/fuel")
				if err = tmpl.ExecuteTemplate(w, "fuel_form", models.FuelFromFuelSheet(csrf.Token(r), fuelSheet)); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = tmpl.ExecuteTemplate(w, "index.html", models.FuelFromFuelSheet(csrf.Token(r), fuelSheet)); err != nil {
					logrus.WithError(err).Error("parsing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		stateService, err := adapters.NewCookieStateService(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			if r.Header.Get("HX-Request") != "true" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("HX-Push-Url", "/export")
			if err = tmpl.ExecuteTemplate(w, "export_form", models.ExportFromExportSheet(csrf.Token(r))); err != nil {
				logrus.WithError(err).Error("executing template")
			}

		case http.MethodGet:
			statsSheet, err := app.Queries.StatsSheet.Handle(r.Context(), stateService)
			if err != nil {
				http.Redirect(w, r, "/fuel", http.StatusSeeOther)
				return
			}

			if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "stats" {
				w.Header().Set("HX-Push-Url", "/stats")
				if err = tmpl.ExecuteTemplate(w, "calculations_form", models.StatsFromStatsSheet(csrf.Token(r), statsSheet)); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = tmpl.ExecuteTemplate(w, "index.html", models.StatsFromStatsSheet(csrf.Token(r), statsSheet)); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`^[A-Za-z0-9 ]+$`)
		if !re.MatchString(r.URL.Query().Get("name")) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		stateService, err := adapters.NewCookieStateService(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pdf, err := app.Queries.PdfExport.Handle(r.Context(), stateService, r.URL.Query().Get("name"))
		if err != nil {
			logrus.WithError(err).Error("creating pdf")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", r.URL.Query().Get("name")))

		_, err = io.Copy(w, pdf)
		if err != nil {
			logrus.WithError(err).Error("writing attachment")
		}
	})

	mux.HandleFunc("/export", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		switch r.Method {
		case http.MethodPost:
			if r.Header.Get("HX-Request") != "true" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = tmpl.ExecuteTemplate(w, "download_redirect.html", models.DownloadRedirectFromStatsSheet(csrf.Token(r), r.Form.Get("reference"))); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		case http.MethodGet:
			if err = tmpl.ExecuteTemplate(w, "index.html", models.ExportFromExportSheet(csrf.Token(r))); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/wind-option", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.Method {
		case http.MethodGet:
			if err = tmpl.ExecuteTemplate(w, "wb_form_wind_option", models.WindOptionsFromRequest(r)); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/fuel-option", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.Method {
		case http.MethodGet:
			if err = tmpl.ExecuteTemplate(w, "fuel_form_max_fuel_option", models.FuelOptionFromRequest(r)); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	CSRF := csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.CookieName("_csrf"))
	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: SecurityHeaders(CSRF(mux)),
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

func SecurityHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Permissions-Policy", "accelerometer=(), autoplay=(), camera=(), cross-origin-isolated=(), display-capture=(), encrypted-media=(), fullscreen=(), geolocation=(), gyroscope=(), keyboard-map=(), magnetometer=(), microphone=(), midi=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
		w.Header().Set("Content-Security-Policy", "script-src 'self' 'unsafe-eval'")
		handler.ServeHTTP(w, r)
	})
}
