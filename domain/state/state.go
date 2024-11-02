package state

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/volume"
	"ppl-calculations/domain/weight"
	"ppl-calculations/domain/wind"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidState = errors.New("invalid state")
)

type Step string

const (
	StepWeight Step = "weight"
	StepFuel   Step = "fuel"
	StepStats  Step = "stats"
	StepExport Step = "export"
)

type State struct {
	// Weight State
	CallSign              *callsign.CallSign       `json:"callSign,omitempty"`
	Pilot                 *weight.Weight           `json:"pilot,omitempty"`
	PilotSeat             *seat.Position           `json:"pilotSeat,omitempty"`
	Passenger             *weight.Weight           `json:"passenger,omitempty"`
	PassengerSeat         *seat.Position           `json:"passengerSeat,omitempty"`
	Baggage               *weight.Weight           `json:"baggage,omitempty"`
	OutsideAirTemperature *temperature.Temperature `json:"outsideAirTemperature,omitempty"`
	PressureAltitude      *pressure.Altitude       `json:"pressureAltitude,omitempty"`
	Wind                  *wind.Wind               `json:"wind,omitempty"`

	// Fuel State
	FuelType          *fuel.Type     `json:"fuelType,omitempty"`
	FuelVolumeType    *volume.Type   `json:"fuelVolumeType,omitempty"`
	Fuel              *fuel.Fuel     `json:"fuel,omitempty"`
	MaxFuel           *bool          `json:"maxFuel,omitempty"`
	TripDuration      *time.Duration `json:"tripDuration,omitempty"`
	AlternateDuration *time.Duration `json:"alternateDuration,omitempty"`
}

type BaseState struct {
	Step string
}

type WeightState struct {
	BaseState
	CallSign              *string
	Pilot                 *string
	PilotSeat             *string
	Passenger             *string
	PassengerSeat         *string
	Baggage               *string
	OutsideAirTemperature *string
	PressureAltitude      *string
	Wind                  *string
	WindDirection         *string
}

type FuelState struct {
	BaseState
	FuelType          string
	FuelVolumeUnit    string
	TripDuration      *string
	AlternateDuration *string
	FuelVolume        *string
	FuelMax           bool
}

type StatsState struct {
	BaseState

	FuelSufficient  bool
	FuelTaxi        string
	FuelTrip        string
	FuelAlternate   string
	FuelContingency string
	FuelReserve     string
	FuelTotal       string
	FuelExtra       string
	FuelExtraAbs    string
}

func (state *State) WriteState(w http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:     "state",
		Value:    state.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	return nil
}

func (state *State) String() string {
	jsonState, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(jsonState)
}

func (state *State) StringPointer(value string) *string {
	return &value
}

func (state *State) WeightState() interface{} {
	is := WeightState{
		BaseState: BaseState{
			Step: string(StepWeight),
		},
	}

	if state.CallSign != nil {
		is.CallSign = state.StringPointer(state.CallSign.String())
	}

	if state.Pilot != nil {
		is.Pilot = state.StringPointer(state.Pilot.String())
	}

	if state.PilotSeat != nil {
		is.PilotSeat = state.StringPointer(state.PilotSeat.String())
	}

	if state.Passenger != nil {
		is.Passenger = state.StringPointer(state.Passenger.String())
	}

	if state.PassengerSeat != nil {
		is.PassengerSeat = state.StringPointer(state.PassengerSeat.String())
	}

	if state.Baggage != nil {
		is.Baggage = state.StringPointer(state.Baggage.String())
	}

	if state.OutsideAirTemperature != nil {
		is.OutsideAirTemperature = state.StringPointer(state.OutsideAirTemperature.String())
	}

	if state.PressureAltitude != nil {
		is.PressureAltitude = state.StringPointer(state.PressureAltitude.String())
	}

	if state.Wind != nil {
		is.Wind = state.StringPointer(state.Wind.Speed.String())
		is.WindDirection = state.StringPointer(state.Wind.Direction.String())
	}

	return is
}

func (state *State) FuelState() interface{} {
	s := FuelState{
		BaseState: BaseState{
			Step: string(StepFuel),
		},
		FuelType:       "mogas",
		FuelVolumeUnit: "liter",
		FuelMax:        false,
	}

	if state.MaxFuel != nil {
		s.FuelMax = *state.MaxFuel
	}

	if state.FuelType != nil {
		s.FuelType = state.FuelType.String()
	}

	if state.FuelVolumeType != nil {
		s.FuelVolumeUnit = state.FuelVolumeType.String()
	}

	if state.Fuel != nil {
		s.FuelVolume = state.StringPointer(fmt.Sprintf("%.1f", state.Fuel.Volume.Amount))
	}

	if state.TripDuration != nil {
		s.TripDuration = state.StringPointer(fmt.Sprintf("%d:%d", int(state.TripDuration.Hours()), int(state.TripDuration.Minutes())%60))
	}

	if state.AlternateDuration != nil {
		s.AlternateDuration = state.StringPointer(fmt.Sprintf("%d:%d", int(state.AlternateDuration.Hours()), int(state.AlternateDuration.Minutes())%60))
	}

	return s
}

func (state *State) parseHHMMToDuration(input string) (time.Duration, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format, expected HH:mm")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour value: %v", err)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute value: %v", err)
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	return duration, nil
}

func (state *State) StatsState() interface{} {
	s := StatsState{
		BaseState: BaseState{
			Step: string(StepStats),
		},
	}

	if state.MaxFuel != nil && *state.MaxFuel {
		planning, err := calculations.NewMaxFuelPlanning(*state.FuelType, *state.TripDuration, *state.AlternateDuration)
		if err != nil {
			panic(err)
		}

		s.FuelSufficient = planning.Sufficient
		s.FuelTaxi = planning.Taxi.Volume.String()
		s.FuelTrip = planning.Trip.Volume.String()
		s.FuelAlternate = planning.Alternate.Volume.String()
		s.FuelContingency = planning.Contingency.Volume.String()
		s.FuelReserve = planning.Reserve.Volume.String()
		s.FuelTotal = planning.Total.Volume.String()
		s.FuelExtra = planning.Extra.Volume.String()
		s.FuelExtraAbs = planning.Extra.Volume.String()

	} else {
		panic("not implemented")
	}

	return s
}

func NewFromStatsRequest(r *http.Request) (*State, error) {
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

func NewFromWeightRequest(r *http.Request) (*State, error) {
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

func NewFromFuelRequest(r *http.Request) (*State, error) {
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
		d, err := s.parseHHMMToDuration(tripDuration)
		if err != nil {
			return nil, err
		}

		s.TripDuration = &d
	}

	if alternateDuration := r.URL.Query().Get("alternate_duration"); alternateDuration != "" {
		d, err := s.parseHHMMToDuration(alternateDuration)
		if err != nil {
			return nil, err
		}

		s.AlternateDuration = &d
	}

	return s, nil
}

func NewFromRequest(r *http.Request) (*State, error) {
	s := &State{}

	if c, err := r.Cookie("state"); err == nil {
		base64DecodedState, err := base64.StdEncoding.DecodeString(c.Value)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(base64DecodedState, &s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
