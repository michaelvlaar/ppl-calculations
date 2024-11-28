package models

import (
	"net/http"
	"ppl-calculations/app/queries"
)

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

type WindOption struct {
	Wind          *string
	WindDirection *string
}

func WindOptionsFromRequest(r *http.Request) interface{} {
	is := WindOption{}

	is.Wind = StringPointer(r.URL.Query().Get("wind"))
	is.WindDirection = StringPointer(r.URL.Query().Get("wind_direction"))

	return is
}

func WeightFromLoadSheet(loadSheet queries.LoadSheetResponse) interface{} {
	is := Weight{
		Base: Base{
			Step: string(StepWeight),
		},
	}

	if loadSheet.CallSign != nil {
		is.CallSign = StringPointer(loadSheet.CallSign.String())
	}

	if loadSheet.Pilot != nil {
		is.Pilot = StringPointer(loadSheet.Pilot.String())
	}

	if loadSheet.PilotSeat != nil {
		is.PilotSeat = StringPointer(loadSheet.PilotSeat.String())
	}

	if loadSheet.Passenger != nil {
		is.Passenger = StringPointer(loadSheet.Passenger.String())
	}

	if loadSheet.PassengerSeat != nil {
		is.PassengerSeat = StringPointer(loadSheet.PassengerSeat.String())
	}

	if loadSheet.Baggage != nil {
		is.Baggage = StringPointer(loadSheet.Baggage.String())
	}

	if loadSheet.OutsideAirTemperature != nil {
		is.OutsideAirTemperature = StringPointer(loadSheet.OutsideAirTemperature.String())
	}

	if loadSheet.PressureAltitude != nil {
		is.PressureAltitude = StringPointer(loadSheet.PressureAltitude.String())
	}

	if loadSheet.Wind != nil {
		is.Wind = StringPointer(loadSheet.Wind.Speed.String())
		is.WindDirection = StringPointer(loadSheet.Wind.Direction.String())
	}

	return is
}
