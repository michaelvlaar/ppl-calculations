package parsing

import (
	"errors"
	"net/http"
	"ppl-calculations/app/commands"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/state"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/volume"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
	"strconv"
)

var (
	ErrInvalidState = errors.New("invalid state")
)

func NewFromRequest(r *http.Request) (*state.State, error) {
	if _, err := r.Cookie("state"); err == nil {
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

func UpdateLoadSheetRequest(r *http.Request) (commands.UpdateLoadSheetRequest, error) {
	req := commands.UpdateLoadSheetRequest{}

	if urlCS := r.URL.Query().Get("callsign"); urlCS != "" {
		cs, err := callsign.New(urlCS)
		if err != nil {
			return req, err
		}

		req.CallSign = &cs
	}

	if urlPilot := r.URL.Query().Get("pilot"); urlPilot != "" {
		pilot, err := weight_balance.NewMassFromString(urlPilot)
		if err != nil {
			return req, err
		}

		req.Pilot = &pilot
	}

	if urlPilotSeat := r.URL.Query().Get("pilot_seat"); urlPilotSeat != "" {
		pilotSeatPosition, err := seat.NewFromString(urlPilotSeat)
		if err != nil {
			return req, err
		}

		req.PilotSeat = &pilotSeatPosition
	}

	if urlPassenger := r.URL.Query().Get("passenger"); urlPassenger != "" {
		passenger, err := weight_balance.NewMassFromString(urlPassenger)
		if err != nil {
			return req, err
		}

		req.Passenger = &passenger
	}

	if urlPassengerSeat := r.URL.Query().Get("passenger_seat"); urlPassengerSeat != "" {
		passengerSeatPosition, err := seat.NewFromString(urlPassengerSeat)
		if err != nil {
			return req, err
		}

		req.PassengerSeat = &passengerSeatPosition
	}

	if urlBaggage := r.URL.Query().Get("baggage"); urlBaggage != "" {
		baggage, err := weight_balance.NewMassFromString(urlBaggage)
		if err != nil {
			return req, err
		}

		req.Baggage = &baggage
	}

	if urlOAT := r.URL.Query().Get("oat"); urlOAT != "" {
		temp, err := temperature.NewFromString(urlOAT)
		if err != nil {
			return req, err
		}

		req.OutsideAirTemperature = &temp
	}

	if urlPA := r.URL.Query().Get("pressure_altitude"); urlPA != "" {
		pa, err := pressure.NewFromString(urlPA)
		if err != nil {
			return req, err
		}

		req.PressureAltitude = &pa
	}

	if urlWind, urlDirection := r.URL.Query().Get("wind"), r.URL.Query().Get("wind_direction"); urlWind != "" && urlDirection != "" {
		sp, err := wind.NewSpeedFromString(urlWind)
		if err != nil {
			return req, err
		}

		d, err := wind.NewDirectionFromString(urlDirection)
		if err != nil {
			return req, err
		}

		w, err := wind.New(d, sp)
		if err != nil {
			return req, err
		}

		req.Wind = &w
	}

	return req, nil
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
		//d, err := parseHHMMToDuration(tripDuration)
		//if err != nil {
		//	return nil, err
		//}

		//s.TripDuration = &d
	}

	if alternateDuration := r.URL.Query().Get("alternate_duration"); alternateDuration != "" {
		//d, err := parseHHMMToDuration(alternateDuration)
		//if err != nil {
		//	return nil, err
		//}

		//s.AlternateDuration = &d
	}

	return s, nil
}
