package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"ppl-calculations/adapters/templator/models"
	"ppl-calculations/adapters/templator/parsing"
)

//go:embed assets/*
var assets embed.FS

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func main() {
	templatesFS, err := fs.Sub(assets, "assets/templates")
	if err != nil {
		log.Fatalf("Fout bij het strippen van templates folder: %v", err)
	}

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"derefString": derefString,
		"mod": func(i, j int) bool {
			return i%j != 0
		},
	}).ParseFS(templatesFS, "*.html")
	if err != nil {
		log.Fatalf("Fout bij het parsen van de templates: %v", err)
	}

	cssFs, err := fs.Sub(assets, "assets/css")
	if err != nil {
		log.Fatalf("Fout bij het parsen van css: %v", err)
	}

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.FS(cssFs))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := parsing.NewFromWeightRequest(r)
		if err != nil {
			_ = tmpl.ExecuteTemplate(w, "index.html", models.WeightFromState(*s))
			return
		}

		_ = parsing.WriteState(s, w)

		if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "weight" {
			w.Header().Set("HX-Push-Url", "/")
			_ = tmpl.ExecuteTemplate(w, "wb_form", models.WeightFromState(*s))
		} else if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Volgende" {
			w.Header().Set("HX-Push-Url", "/fuel")
			_ = tmpl.ExecuteTemplate(w, "fuel_form", models.FuelFromState(*s))
		} else {
			_ = tmpl.ExecuteTemplate(w, "index.html", models.WeightFromState(*s))
		}

	})

	http.HandleFunc("/fuel", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := parsing.NewFromFuelRequest(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		_ = parsing.WriteState(s, w)

		if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Vorige" {
			w.Header().Set("HX-Push-Url", "/")
			_ = tmpl.ExecuteTemplate(w, "wb_form", models.WeightFromState(*s))
		} else if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Volgende" {
			w.Header().Set("HX-Push-Url", "/stats")
			_ = tmpl.ExecuteTemplate(w, "calculations_form", models.StatsFromState(*s))
		} else {
			_ = tmpl.ExecuteTemplate(w, "index.html", models.FuelFromState(*s))
		}
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := parsing.NewFromStatsRequest(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		_ = parsing.WriteState(s, w)

		if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Vorige" {
			w.Header().Set("HX-Push-Url", "/fuel")
			_ = tmpl.ExecuteTemplate(w, "fuel_form", models.FuelFromState(*s))
		} else if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("submit") == "Volgende" {
			w.Header().Set("HX-Push-Url", "/export")
			_ = tmpl.ExecuteTemplate(w, "calculations_form", models.StatsFromState(*s))
		} else {
			_ = tmpl.ExecuteTemplate(w, "index.html", models.StatsFromState(*s))
		}
	})

	http.HandleFunc("/wind-option", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := parsing.NewFromWeightRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "wb_form_wind_option", models.WeightFromState(*s))
	})

	http.HandleFunc("/fuel-option", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		s, err := parsing.NewFromFuelRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "fuel_form_max_fuel_option", models.FuelFromState(*s))
	})

	fmt.Println("Server gestart op :80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		return
	}
}
