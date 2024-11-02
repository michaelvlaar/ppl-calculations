package models

import "ppl-calculations/domain/state"

type Weight struct {
	Base

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

func WeightFromState(s state.State) interface{} {
	is := Weight{
		Base: Base{
			Step: string(StepWeight),
		},
	}

	if s.CallSign != nil {
		is.CallSign = StringPointer(s.CallSign.String())
	}

	if s.Pilot != nil {
		is.Pilot = StringPointer(s.Pilot.String())
	}

	if s.PilotSeat != nil {
		is.PilotSeat = StringPointer(s.PilotSeat.String())
	}

	if s.Passenger != nil {
		is.Passenger = StringPointer(s.Passenger.String())
	}

	if s.PassengerSeat != nil {
		is.PassengerSeat = StringPointer(s.PassengerSeat.String())
	}

	if s.Baggage != nil {
		is.Baggage = StringPointer(s.Baggage.String())
	}

	if s.OutsideAirTemperature != nil {
		is.OutsideAirTemperature = StringPointer(s.OutsideAirTemperature.String())
	}

	if s.PressureAltitude != nil {
		is.PressureAltitude = StringPointer(s.PressureAltitude.String())
	}

	if s.Wind != nil {
		is.Wind = StringPointer(s.Wind.Speed.String())
		is.WindDirection = StringPointer(s.Wind.Direction.String())
	}

	return is
}
