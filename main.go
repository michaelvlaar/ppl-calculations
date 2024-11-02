package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"ppl-calculations/domain/state"
)

//go:embed templates/*
var content embed.FS

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func main() {
	templatesFS, err := fs.Sub(content, "templates")
	if err != nil {
		log.Fatalf("Fout bij het strippen van templates folder: %v", err)
	}

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"derefString": derefString,
	}).ParseFS(templatesFS, "*.html")
	if err != nil {
		log.Fatalf("Fout bij het parsen van de templates: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := state.NewFromWeightRequest(r)
		if err != nil {
			_ = tmpl.ExecuteTemplate(w, "index.html", s.WeightState())
			return
		}

		_ = s.WriteState(w)

		if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "weight" {
			w.Header().Set("HX-Push-Url", "/")
			_ = tmpl.ExecuteTemplate(w, "wb_form", s.WeightState())
		} else if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Volgende" {
			w.Header().Set("HX-Push-Url", "/fuel")
			_ = tmpl.ExecuteTemplate(w, "fuel_form", s.FuelState())
		} else {
			err = tmpl.ExecuteTemplate(w, "index.html", s.WeightState())
			if err != nil {
				panic(err)
			}
		}

	})

	http.HandleFunc("/fuel", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := state.NewFromFuelRequest(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		_ = s.WriteState(w)

		if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Vorige" {
			w.Header().Set("HX-Push-Url", "/")
			_ = tmpl.ExecuteTemplate(w, "wb_form", s.WeightState())
		} else if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Volgende" {
			w.Header().Set("HX-Push-Url", "/stats")
			_ = tmpl.ExecuteTemplate(w, "calculations_form", s.StatsState())
		} else {
			_ = tmpl.ExecuteTemplate(w, "index.html", s.FuelState())
		}
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := state.NewFromStatsRequest(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		_ = s.WriteState(w)

		if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Vorige" {
			w.Header().Set("HX-Push-Url", "/fuel")
			_ = tmpl.ExecuteTemplate(w, "fuel_form", s.FuelState())
		} else if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Volgende" {
			w.Header().Set("HX-Push-Url", "/export")
			_ = tmpl.ExecuteTemplate(w, "calculations_form", s.StatsState())
		} else {
			_ = tmpl.ExecuteTemplate(w, "index.html", s.StatsState())
		}
	})

	http.HandleFunc("/wind-option", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := state.NewFromWeightRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "wb_form_wind_option", s.WeightState())
	})

	http.HandleFunc("/fuel-option", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := state.NewFromFuelRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "fuel_form_max_fuel_option", s.FuelState())
	})

	fmt.Println("Server gestart op :80")
	http.ListenAndServe(":80", nil)
}
