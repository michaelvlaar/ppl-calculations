package main

import (
	"embed"
	"fmt"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"html/template"
	"image/color"
	"io/fs"
	"log"
	"net/http"
	"ppl-calculations/adapters/templator/models"
	"ppl-calculations/adapters/templator/parsing"
	"ppl-calculations/domain/calculations"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

//go:embed assets/*
var assets embed.FS

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func TickGenerator(min, max float64) plot.Ticker {
	return plot.ConstantTicks(func() []plot.Tick {
		var ticks []plot.Tick
		step := 20.0 // Ticks per 10 units
		for val := min; val <= max; val += step {
			ticks = append(ticks, plot.Tick{
				Value: val,
				Label: fmt.Sprintf("%.0f", val),
			})
		}
		return ticks
	}())
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

	http.HandleFunc("/aquila-wb", func(w http.ResponseWriter, r *http.Request) {
		urlTakeOffMass := r.URL.Query().Get("takeoff-mass")
		urlTakeOffMassMoment := r.URL.Query().Get("takeoff-mass-moment")

		urlLandingMass := r.URL.Query().Get("landing-mass")
		urlLandingMassMoment := r.URL.Query().Get("landing-mass-moment")

		if urlTakeOffMass == "" || urlTakeOffMassMoment == "" || urlLandingMass == "" || urlLandingMassMoment == "" {
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

		p := plot.New()
		p.Title.Text = "Gewicht en Balans (PHDHA)"
		p.X.Label.Text = "Mass Moment [kg m]"
		p.Y.Label.Text = "Mass [kg]"

		p.X.Min = 230
		p.X.Max = 430
		p.Y.Min = 550
		p.Y.Max = 770

		p.X.Tick.Marker = TickGenerator(p.X.Min, p.X.Max)
		p.Y.Tick.Marker = TickGenerator(p.Y.Min, p.Y.Max)

		grid := plotter.NewGrid()
		p.Add(grid)

		polygonData := plotter.XYs{
			{X: calculations.AquilaMinWeight * calculations.AquilaForwardCgLimit, Y: calculations.AquilaMinWeight},
			{X: calculations.AquilaMinWeight * calculations.AquilaBackwardCgLimit, Y: calculations.AquilaMinWeight},
			{X: calculations.AquilaMTOW * calculations.AquilaBackwardCgLimit, Y: calculations.AquilaMTOW},
			{X: calculations.AquilaMTOW * calculations.AquilaForwardCgLimit, Y: calculations.AquilaMTOW},
		}

		polygon, err := plotter.NewPolygon(polygonData)
		if err != nil {
			log.Fatalf("Error creating polygon: %v", err)
		}

		polygon.LineStyle.Color = color.RGBA{R: 255, A: 0}
		polygon.Color = color.RGBA{R: 255, A: 128}
		p.Add(polygon)

		points := plotter.XYs{
			{X: takeOffMassMoment, Y: takeOffMass},
		}

		points2 := plotter.XYs{
			{X: landingMassMoment, Y: landingMass},
		}

		scatter, err := plotter.NewScatter(points)
		if err != nil {
			log.Fatalf("Error creating scatter: %v", err)
		}

		scatter.GlyphStyle.Shape = draw.CircleGlyph{}
		scatter.GlyphStyle.Color = color.RGBA{G: 255, A: 255}
		scatter.GlyphStyle.Radius = vg.Points(5)

		p.Add(scatter)

		scatter2, err := plotter.NewScatter(points2)
		if err != nil {
			log.Fatalf("Error creating scatter: %v", err)
		}

		scatter2.GlyphStyle.Shape = draw.CircleGlyph{}
		scatter2.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
		scatter2.GlyphStyle.Radius = vg.Points(5)

		p.Add(scatter2)

		p.Legend.Top = false
		p.Legend.XOffs = vg.Centimeter * -1.5
		p.Legend.YOffs = vg.Centimeter * 1.5
		p.Legend.Padding = vg.Points(5)

		p.Legend.Add("CG Envelope", polygon)
		p.Legend.Add("Take-off Point", scatter)
		p.Legend.Add("Landing Point", scatter2)

		c := canvas.New(150.0, 150.0)
		gonumCanvas := renderers.NewGonumPlot(c)

		p.Draw(gonumCanvas)

		w.Header().Set("content-type", "image/svg+xml")
		err = c.Write(w, renderers.SVG())
		if err != nil {
			panic(err)
		}
	})

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
