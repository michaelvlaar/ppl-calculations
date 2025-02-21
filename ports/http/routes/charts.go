package routes

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"ppl-calculations/app"
	"ppl-calculations/app/queries"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
	"ppl-calculations/ports/http/middleware"
	"strconv"
)

func RegisterChartRoutes(mux *http.ServeMux, app app.Application) {
	mux.HandleFunc("/aquila-wb", middleware.CacheHeadersFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	mux.HandleFunc("/aquila-tdr", middleware.CacheHeadersFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	mux.HandleFunc("/aquila-ldr", middleware.CacheHeadersFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))

}
