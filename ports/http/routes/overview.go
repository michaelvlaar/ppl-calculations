package routes

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"ppl-calculations/adapters"
	"ppl-calculations/app"
	"ppl-calculations/app/commands"
	"ppl-calculations/domain/export"
	"ppl-calculations/ports/http/middleware"
	"ppl-calculations/ports/templates"
	"ppl-calculations/ports/templates/models"
)

func RegisterOverviewRoutes(mux *http.ServeMux, app app.Application) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		stateService, err := adapters.NewCookieStateService(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodDelete:
			if r.Header.Get("HX-Request") != "true" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			id, err := export.NewIDFromString(r.Form.Get("id"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = app.Commands.DeleteExportSheet.Handle(r.Context(), stateService, commands.DeleteExportSheetRequest{
				ID: id,
			}); err != nil {
				logrus.WithError(err).Error("delete export")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ex, err := app.Queries.Exports.Handle(r.Context(), stateService)
			if len(ex.Exports) == 0 {
				if err = templates.OverviewNoExports().Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			}
			return
		case http.MethodGet:
			respEx, err := app.Queries.Exports.Handle(r.Context(), stateService)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			tmpl := templates.Overview(models.OverviewFromExports(respEx.Exports))
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Push-Url", "/")
				if err = tmpl.Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			} else {
				if err = templates.Index(tmpl).Render(r.Context(), w); err != nil {
					logrus.WithError(err).Error("executing template")
				}
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/view", middleware.CacheHeadersFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		switch r.Method {
		case http.MethodGet:
			d := r.URL.Query().Get("d")
			b, err := base64.URLEncoding.DecodeString(d)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			gzReader, err := gzip.NewReader(bytes.NewReader(b))
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			gzBytes, err := io.ReadAll(gzReader)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			err = gzReader.Close()
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			buf := bytes.NewBuffer(gzBytes)
			dec := gob.NewDecoder(buf)

			var e export.Export
			if err := dec.Decode(&e); err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			statsSheet, err := app.Queries.StatsSheet.HandleExport(r.Context(), e)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			templ := templates.CalculationsView(models.ViewFromExport(statsSheet, e))
			if r.Header.Get("HX-Request") == "true" {
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
	}))

	mux.HandleFunc("/download", middleware.CacheHeadersFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		switch r.Method {
		case http.MethodGet:
			d := r.URL.Query().Get("d")
			b, err := base64.URLEncoding.DecodeString(d)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			gzReader, err := gzip.NewReader(bytes.NewReader(b))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			gzBytes, err := io.ReadAll(gzReader)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = gzReader.Close()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			buf := bytes.NewBuffer(gzBytes)
			dec := gob.NewDecoder(buf)

			var e export.Export
			if err := dec.Decode(&e); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			pdf, err := app.Queries.PdfExport.Handle(r.Context(), e)
			if err != nil {
				logrus.WithError(err).Error("creating pdf")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", e.Name.String()))

			_, err = io.Copy(w, pdf)
			if err != nil {
				logrus.WithError(err).Error("writing attachment")
			}
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}))
}
