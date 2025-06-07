package routes

import (
	"github.com/michaelvlaar/ppl-calculations/app"
	"github.com/michaelvlaar/ppl-calculations/ports/templates"
	"github.com/michaelvlaar/ppl-calculations/ports/templates/models"
	"github.com/michaelvlaar/ppl-calculations/ports/templates/parsing"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RegisterCalculationRoutes(mux *http.ServeMux, app app.Application) {
	mux.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
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

			updateLoadSheetRequest, err := parsing.UpdateLoadSheetRequest(r)
			if err != nil {
				logrus.WithError(err).Error("creating update loadsheet request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = app.Commands.UpdateLoadSheet.Handle(r.Context(), updateLoadSheetRequest); err != nil {
				logrus.WithError(err).Error("update loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fuelSheet, err := app.Queries.FuelSheet.Handle(r.Context())
			if err != nil {
				logrus.WithError(err).Error("reading fuelsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Push-Url", "/fuel")

			if err = templates.Fuel(models.FuelFromFuelSheet(fuelSheet)).Render(r.Context(), w); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		case http.MethodGet:
			loadSheet, err := app.Queries.LoadSheet.Handle(r.Context())
			if err != nil {
				logrus.WithError(err).Error("reading loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			templ := templates.WeightAndBalanceForm(models.WeightFromLoadSheet(loadSheet))
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Push-Url", "/load")
				if err = templ.Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = templates.Index(templ).Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/fuel", func(w http.ResponseWriter, r *http.Request) {
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

			updateFuelSheetRequest, err := parsing.UpdateFuelSheetRequest(r)
			if err != nil {
				logrus.WithError(err).Error("creating update fuelsheet request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = app.Commands.UpdateFuelSheet.Handle(r.Context(), updateFuelSheetRequest); err != nil {
				logrus.WithError(err).Error("update loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			statsSheet, err := app.Queries.StatsSheet.Handle(r.Context())
			if err != nil {
				logrus.WithError(err).Error("update loadsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Push-Url", "/stats")
			if err = templates.CalculationsForm(models.StatsFromStatsSheet(statsSheet)).Render(r.Context(), w); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		case http.MethodGet:
			fuelSheet, err := app.Queries.FuelSheet.Handle(r.Context())
			if err != nil {
				http.Redirect(w, r, "/load", http.StatusSeeOther)
				return
			}

			templ := templates.Fuel(models.FuelFromFuelSheet(fuelSheet))
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Push-Url", "/fuel")
				if err = templ.Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = templates.Index(templ).Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("parsing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
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

			sheet, err := app.Queries.ExportSheet.Handle(r.Context())
			if err != nil {
				logrus.WithError(err).Error("loading exportsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Push-Url", "/export")
			if err = templates.Export(models.ExportFromExportSheet(sheet)).Render(r.Context(), w); err != nil {
				logrus.WithError(err).Error("executing template")
			}

		case http.MethodGet:
			statsSheet, err := app.Queries.StatsSheet.Handle(r.Context())
			if err != nil {
				http.Redirect(w, r, "/fuel", http.StatusSeeOther)
				return
			}

			templ := templates.CalculationsForm(models.StatsFromStatsSheet(statsSheet))
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Push-Url", "/stats")
				if err = templ.Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = templates.Index(templ).Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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

			updateExportSheetRequest, err := parsing.UpdateExportSheetRequest(r)
			if err != nil {
				logrus.WithError(err).Error("creating update fuelsheet request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = app.Commands.UpdateExportSheet.Handle(r.Context(), updateExportSheetRequest); err != nil {
				logrus.WithError(err).Error("update exportsheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err = app.Commands.ClearSheet.Handle(r.Context()); err != nil {
				logrus.WithError(err).Error("clear sheet")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			respEx, err := app.Queries.Exports.Handle(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("HX-Push-Url", "/")
			if err = templates.Overview(models.OverviewFromExports(respEx.Exports)).Render(r.Context(), w); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		case http.MethodGet:
			sheet, err := app.Queries.ExportSheet.Handle(r.Context())
			if err != nil {
				logrus.WithError(err).Error("creating pdf")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			templ := templates.Export(models.ExportFromExportSheet(sheet))
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Push-Url", "/export")
				if err = templ.Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = templates.Index(templ).Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/wind-option", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.Method {
		case http.MethodGet:
			if err := templates.WeightAndBalanceWindOption(models.WindOptionsFromRequest(r)).Render(r.Context(), w); err != nil {
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
			if err := templates.FuelMaxFuel(models.FuelOptionFromRequest(r)).Render(r.Context(), w); err != nil {
				logrus.WithError(err).Error("executing template")
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}
