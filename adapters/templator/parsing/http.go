package parsing

import (
	"errors"
	"net/http"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/state"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/volume"
	"ppl-calculations/domain/weight"
	"ppl-calculations/domain/wind"
	"strconv"
	"time"
)

var (
	ErrInvalidState = errors.New("invalid state")
)

func WriteState(s *state.State, w http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:     "state",
		Value:    s.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	return nil
}

func NewFromRequest(r *http.Request) (*state.State, error) {
	if c, err := r.Cookie("state"); err == nil {
		return state.NewFromString(c.Value)
	}

	return state.MustNew(), nil
}
func NewFromStatsRequest(r *http.Request) (*state.State, error) {
	s, err := NewFromRequest(r)
	if err != nil {
		return nil, err
	}

	// Verify if state is present
	// TODO: improve this check, how to determine if the state is valid?
	if s.MaxFuel == nil {
		return s, ErrInvalidState
	}

	return s, nil
}

func NewFromWeightRequest(r *http.Request) (*state.State, error) {
	s, err := NewFromRequest(r)
	if err != nil {
		return nil, err
	}

	if urlCS := r.URL.Query().Get("callsign"); urlCS != "" {
		cs, err := callsign.New(urlCS)
		if err != nil {
			return nil, err
		}

		s.CallSign = &cs
	}

	if urlPilot := r.URL.Query().Get("pilot"); urlPilot != "" {
		pilot, err := weight.NewFromString(urlPilot)
		if err != nil {
			return nil, err
		}

		s.Pilot = &pilot
	}

	if urlPilotSeat := r.URL.Query().Get("pilot_seat"); urlPilotSeat != "" {
		pilotSeatPosition, err := seat.NewFromString(urlPilotSeat)
		if err != nil {
			return nil, err
		}

		s.PilotSeat = &pilotSeatPosition
	}

	if urlPassenger := r.URL.Query().Get("passenger"); urlPassenger != "" {
		passenger, err := weight.NewFromString(urlPassenger)
		if err != nil {
			return nil, err
		}

		s.Passenger = &passenger
	}

	if urlPassengerSeat := r.URL.Query().Get("passenger_seat"); urlPassengerSeat != "" {
		passengerSeatPosition, err := seat.NewFromString(urlPassengerSeat)
		if err != nil {
			return nil, err
		}

		s.PassengerSeat = &passengerSeatPosition
	}

	if urlBaggage := r.URL.Query().Get("baggage"); urlBaggage != "" {
		baggage, err := weight.NewFromString(urlBaggage)
		if err != nil {
			return nil, err
		}

		s.Baggage = &baggage
	}

	if urlOAT := r.URL.Query().Get("oat"); urlOAT != "" {
		temp, err := temperature.NewFromString(urlOAT)
		if err != nil {
			return nil, err
		}

		s.OutsideAirTemperature = &temp
	}

	if urlPA := r.URL.Query().Get("pressure_altitude"); urlPA != "" {
		pa, err := pressure.NewFromString(urlPA)
		if err != nil {
			return nil, err
		}

		s.PressureAltitude = &pa
	}

	if urlWind, urlDirection := r.URL.Query().Get("wind"), r.URL.Query().Get("wind_direction"); urlWind != "" && urlDirection != "" {
		sp, err := wind.NewSpeedFromString(urlWind)
		if err != nil {
			return nil, err
		}

		d, err := wind.NewDirectionFromString(urlDirection)
		if err != nil {
			return nil, err
		}

		w, err := wind.New(d, sp)
		if err != nil {
			return nil, err
		}

		s.Wind = &w
	}

	return s, nil
}

func NewFromFuelRequest(r *http.Request) (*state.State, error) {
	s, err := NewFromRequest(r)
	if err != nil {
		return nil, err
	}

	// Verify if state is present
	// TODO: improve this check, how to determine if the state is valid?
	if s.CallSign == nil {
		return s, ErrInvalidState
	}

	maxFuel := r.URL.Query().Get("fuel_max") == "max"
	s.MaxFuel = &maxFuel

	if urlFuelType := r.URL.Query().Get("fuel_type"); urlFuelType != "" {
		fuelType, err := fuel.NewTypeFromString(r.URL.Query().Get("fuel_type"))
		if err != nil {
			return nil, err
		}

		s.FuelType = &fuelType
	}

	if urlFuelUnit := r.URL.Query().Get("fuel_unit"); urlFuelUnit != "" {
		fuelUnit, err := volume.NewTypeFromString(r.URL.Query().Get("fuel_unit"))
		if err != nil {
			return nil, err
		}

		s.FuelVolumeType = &fuelUnit
	}

	if fuelVol := r.URL.Query().Get("fuel_volume"); fuelVol != "" && s.FuelType != nil && s.FuelVolumeType != nil && s.MaxFuel != nil && *s.MaxFuel == false {
		vol, err := strconv.ParseFloat(fuelVol, 64)
		if err != nil {
			return nil, err
		}

		v, err := volume.New(vol, *s.FuelVolumeType)
		if err != nil {
			return nil, err
		}

		f, err := fuel.New(v, *s.FuelType)
		if err != nil {
			return nil, err
		}

		s.Fuel = &f
	}

	if tripDuration := r.URL.Query().Get("trip_duration"); tripDuration != "" {
		d, err := parseHHMMToDuration(tripDuration)
		if err != nil {
			return nil, err
		}

		s.TripDuration = &d
	}

	if alternateDuration := r.URL.Query().Get("alternate_duration"); alternateDuration != "" {
		d, err := parseHHMMToDuration(alternateDuration)
		if err != nil {
			return nil, err
		}

		s.AlternateDuration = &d
	}

	return s, nil
}
