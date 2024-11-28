package queries

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
	"text/template"
)

type LdrChartHandler struct {
	chart bytes.Buffer
}

func NewLdrChartHandler(chart bytes.Buffer) LdrChartHandler {
	return LdrChartHandler{
		chart: chart,
	}
}

type LdrChartRequest struct {
	OAT              temperature.Temperature
	PressureAltitude pressure.Altitude
	Tow              *weight_balance.Mass
	Wind             *wind.Wind
}

func interpolateYValues(pressureAltitude float64, oatY [][2]interface{}, yBracket int) (float64, float64, error) {
	var yLow, yHigh []float64
	var altLow, altHigh float64

	switch {
	case pressureAltitude <= 2000:
		altLow, yLow = oatY[0][0].(float64), oatY[0][1].([]float64)
		altHigh, yHigh = oatY[1][0].(float64), oatY[1][1].([]float64)
	case pressureAltitude <= 4000:
		altLow, yLow = oatY[1][0].(float64), oatY[1][1].([]float64)
		altHigh, yHigh = oatY[2][0].(float64), oatY[2][1].([]float64)
	case pressureAltitude <= 6000:
		altLow, yLow = oatY[2][0].(float64), oatY[2][1].([]float64)
		altHigh, yHigh = oatY[3][0].(float64), oatY[3][1].([]float64)
	case pressureAltitude <= 8000:
		altLow, yLow = oatY[3][0].(float64), oatY[3][1].([]float64)
		altHigh, yHigh = oatY[4][0].(float64), oatY[4][1].([]float64)
	default:
		return 0, 0, errors.New("pressure altitude out of range")
	}

	yFactor := (pressureAltitude - altLow) / (altHigh - altLow)

	y1 := interpolate(yLow[yBracket], yHigh[yBracket], yFactor)
	y2 := interpolate(
		yLow[min(yBracket+1, len(yLow)-1)],
		yHigh[min(yBracket+1, len(yHigh)-1)],
		yFactor,
	)

	return y1, y2, nil
}

func interpolateTomY(tom [2][4]float64, mtow, interpolatedY, pOatY, yOffset float64) float64 {
	tomYLow := interpolate(tom[0][2], tom[0][3], (750.0-mtow)/(tom[0][1]-tom[0][0]))
	tomYHigh := interpolate(tom[1][2], tom[1][3], (750.0-mtow)/(tom[1][1]-tom[1][0]))
	return interpolate(tomYLow, tomYHigh, (interpolatedY+(pOatY*yOffset)-tom[0][2])/(tom[1][2]-tom[0][2]))
}

func interpolate(start, end, factor float64) float64 {
	return start + (end-start)*factor
}

type LdrData struct {
	OatX string
	OatY string
	LmX  string
	LmY  string
}

func (h LdrChartHandler) Handle(_ context.Context, request LdrChartRequest) (io.Reader, error) {
	oatXStart := 562.923177
	oatXEnd := 1870.93099
	oatXUnits := 70.0

	oatY := [][2]interface{}{
		{0.0, []float64{1902.34, 1948.34, 1994.34, 2042.32, 2090.33, 2136.33, 2184.34, 2234.34}},
		{2000.0, []float64{2002.34, 2054.33, 2104.33, 2158.33, 2210.32, 2262.34, 2316.34, 2370.34}},
		{4000.0, []float64{2114.32, 2172.33, 2228.32, 2286.33, 2344.34, 2404.33, 2462.34, 2522.33}},
		{6000.0, []float64{2242.32, 2304.33, 2368.33, 2432.32, 2498.34, 2562.33, 2628.32, 2694.34}},
		{8000.0, []float64{2384.34, 2454.33, 2526.33, 2598.34, 2670.34, 2742.32, 2814.33, 2888.31}},
	}

	yBracket := int((request.OAT.Float64() + 30) / 10)
	yInterpolated1, yInterpolated2, err := interpolateYValues(request.PressureAltitude.Float64(), oatY, yBracket)

	if err != nil {
		return nil, err
	}

	pOatX := (oatXEnd - oatXStart) / oatXUnits
	yOffset := math.Mod(request.OAT.Float64()+30.0, 10)
	pOatY := (yInterpolated2 - yInterpolated1) / 10

	log.Print(oatXStart+pOatX*(request.OAT.Float64()+30), yInterpolated1+pOatY*yOffset)

	tomXStart := 2077.115885
	tomXEnd := 3263.216146
	tomUnits := 750.0 - 550.0
	tomXOffset := (750.0-request.Tow.Kilo())*(tomXEnd-tomXStart)/tomUnits + tomXStart

	tom := [2][4]float64{
		{0.0, 200.0, 1906.05, 1796.06},
		{0.0, 200.0, 2002.08, 1882.06},
	}

	tomYPos := interpolateTomY(tom, request.Tow.Kilo(), yInterpolated1, pOatY, yOffset)

	//windXStart := 3439.388021
	//windXEnd := 4933.561198
	//windXPos, windYPos := calculateWindPosition(windXStart, windXEnd, wind, tomYPos, tomXOffset)
	//
	//obs := [2][2]float64{
	//	{1467.55, 1171.48},
	//	{1631.61, 1241.50},
	//}
	//obsYPos := calculateObstacleY(windYPos, obs)

	//return [9]float64{oatXStart + pOatX*(oat+30), yInterpolated1 + pOatY*yOffset, tomXOffset, tomYPos, windXPos, windYPos, obsYPos, 0, 0}, nil

	tmpl, err := template.New("svgTemplate").Parse(h.chart.String())
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, LdrData{
		OatX: fmt.Sprintf("%.5f", oatXStart+pOatX*(request.OAT.Float64()+30.0)),
		OatY: fmt.Sprintf("%.5f", yInterpolated1+pOatY*yOffset),
		LmX:  fmt.Sprintf("%.5f", tomXOffset),
		LmY:  fmt.Sprintf("%.5f", tomYPos),
	})

	if err != nil {
		return nil, err
	}

	return &output, nil
}
