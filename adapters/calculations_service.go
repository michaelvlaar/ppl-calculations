package adapters

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"io"
	"io/fs"
	"math"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/fuel"
	"ppl-calculations/domain/pressure"
	"ppl-calculations/domain/seat"
	"ppl-calculations/domain/temperature"
	"ppl-calculations/domain/weight_balance"
	"ppl-calculations/domain/wind"
	"text/template"
	"time"
)

type CalculationsService struct {
	tdrChart     bytes.Buffer
	ldrChart     bytes.Buffer
	imageService *ImageService
}

func MustNewCalculationsService(assets fs.FS, imageService *ImageService) *CalculationsService {
	ldrFile, err := assets.Open("ldr.svg")
	if err != nil {
		logrus.WithError(err).Fatal("ldr.svg not present")
	}

	var ldrBuf bytes.Buffer
	_, err = io.Copy(&ldrBuf, ldrFile)
	if err != nil {
		logrus.WithError(err).Fatal("ldr.svg not loaded")
	}

	err = ldrFile.Close()
	if err != nil {
		logrus.WithError(err).Fatal("ldr.svg not closed")
	}

	tdrFile, err := assets.Open("tdr.svg")
	if err != nil {
		logrus.WithError(err).Fatal("tdr.svg not present")
	}

	var tdrBuf bytes.Buffer
	_, err = io.Copy(&tdrBuf, tdrFile)
	if err != nil {
		logrus.WithError(err).Fatal("tdr.svg not loaded")
	}

	err = tdrFile.Close()
	if err != nil {
		logrus.WithError(err).Fatal("ldr.svg not closed")
	}

	return &CalculationsService{
		tdrChart:     tdrBuf,
		ldrChart:     ldrBuf,
		imageService: imageService,
	}
}

func (s CalculationsService) interpolate(start, end, factor float64) float64 {
	return start + (end-start)*factor
}

func (s CalculationsService) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s CalculationsService) calculateWindPositionTOD(windXStart, windXEnd, wind, tomYPos, tomXOffset float64) (float64, float64) {
	windXPos := tomXOffset
	windYPos := tomYPos

	if wind != 0.0 {
		windUnits := 20.0
		windXOffset := math.Abs(wind)*((windXEnd-windXStart)/windUnits) + windXStart
		windOffset := wind

		initialFactor := [2][4]float64{
			{0.0, 10.0, 1389.84375, 1303.841146},
			{0.0, 10.0, 1655.891927, 1507.877604},
		}

		var windD [2][4]float64

		if wind >= 0.0 && wind <= 10.0 {
			windD = initialFactor
		} else if wind > 10.0 && wind <= 15.0 {
			windOffset = math.Mod(wind, 10.0)
			windD = [2][4]float64{
				{0.0, 5.0, 1303.841146, 1269.856771},
				{0.0, 5.0, 1507.877604, 1449.869792},
			}
		} else if wind >= -10.0 && wind < 0.0 {
			if tomYPos <= 1640.891927 {
				windD = [2][4]float64{
					{0.0, 10.0, 1389.84375, 1525.84375},
					{0.0, 10.0, 1640.891927, 1867.0},
				}
			} else if tomYPos <= 1958.915365 {
				windD = [2][4]float64{
					{0.0, 10.0, 1640.891927, 1867.0},
					{0.0, 10.0, 1958.915365, 2300.0},
				}
			} else {
				windD = [2][4]float64{
					{0.0, 10.0, 1958.915365, 2300.0},
					{0.0, 10.0, 2262.979167, 2710.0},
				}
			}
		} else {
			windOffset = math.Mod(wind, 15.0)
			windD = [2][4]float64{
				{0.0, 5.0, 1269.856771, 1243.847656},
				{0.0, 5.0, 1449.869792, 1407.845052},
			}
		}

		windLow := windD[0][2] - ((windD[0][2]-windD[0][3])/(windD[0][1]-windD[0][0]))*math.Abs(windOffset)
		windHigh := windD[1][2] - ((windD[1][2]-windD[1][3])/(windD[1][1]-windD[1][0]))*math.Abs(windOffset)

		var windFactor float64
		if wind >= 0.0 {
			windFactor = (tomYPos - initialFactor[0][2]) / (initialFactor[1][2] - initialFactor[0][2])
		} else {
			windFactor = (tomYPos - windD[0][2]) / (windD[1][2] - windD[0][2])
		}

		windYPos = (windHigh-windLow)*windFactor + windLow
		windXPos = math.Abs(windXOffset)
	}

	return windXPos, windYPos
}

func (s CalculationsService) interpolateObstacleY(windYPos float64, obsData [2][2]float64) float64 {
	obsFactor := (windYPos - obsData[0][0]) / (obsData[1][0] - obsData[0][0])
	return obsData[0][1] + obsFactor*(obsData[1][1]-obsData[0][1])
}

func (s CalculationsService) findNextValue(values []float64, target, defaultValue float64) float64 {
	for _, v := range values {
		if v >= target {
			return v
		}
	}
	return defaultValue
}

type TodData struct {
	OatX  string
	OatY  string
	TomX  string
	TomY  string
	WindX string
	WindY string
	ObsY  string
}

type LdrData struct {
	OatX  string
	OatY  string
	LmX   string
	LmY   string
	WindX string
	WindY string
	ObY   string
}

func (s CalculationsService) calculateObstacleY(windYPos float64, obs [2][2]float64) float64 {
	obsFactor := (windYPos - obs[0][0]) / (obs[1][0] - obs[0][0])
	return obs[0][1] + obsFactor*(obs[1][1]-obs[0][1])
}

func (s CalculationsService) interpolateYValues(pressureAltitude float64, oatY [][2]interface{}, yBracket int) (float64, float64, error) {
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

	y1 := s.interpolate(yLow[yBracket], yHigh[yBracket], yFactor)
	y2 := s.interpolate(
		yLow[min(yBracket+1, len(yLow)-1)],
		yHigh[min(yBracket+1, len(yHigh)-1)],
		yFactor,
	)

	return y1, y2, nil
}

func (s CalculationsService) interpolateTomY(tom [2][4]float64, mtow, interpolatedY, pOatY, yOffset float64) float64 {
	tomLowFactor := (750.0 - mtow) / (tom[0][1] - tom[0][0])
	tomHighFactor := (750.0 - mtow) / (tom[1][1] - tom[1][0])

	tomYLow := s.interpolate(tom[0][2], tom[0][3], tomLowFactor)
	tomYHigh := s.interpolate(tom[1][2], tom[1][3], tomHighFactor)

	yFactor := (interpolatedY + (pOatY * yOffset) - tom[0][2]) / (tom[1][2] - tom[0][2])
	return s.interpolate(tomYLow, tomYHigh, yFactor)
}

func (s CalculationsService) calculateWindPosition(windXStart, windXEnd, wind, tomYPos, tomXOffset float64) (float64, float64) {
	windXPos := tomXOffset
	windYPos := tomYPos

	if wind != 0.0 {
		windUnits := 20.0
		windXOffset := math.Abs(wind)*((windXEnd-windXStart)/windUnits) + windXStart
		windOffset := wind

		initialFactor := [2][4]float64{
			{0.0, 10.0, 1787.923177, 1599.902344},
			{0.0, 10.0, 2173.958333, 1897.916667},
		}

		var windD [2][4]float64

		if wind >= 0.0 && wind <= 10.0 {
			windD = initialFactor
		} else if wind > 10.0 && wind <= 15.0 {
			windOffset = math.Mod(wind, 10.0)
			windD = [2][4]float64{
				{0.0, 5.0, 1599.902344, 1527.864583},
				{0.0, 5.0, 1897.916667, 1791.894531},
			}
		} else if wind >= -10.0 && wind < 0.0 {
			if tomYPos <= 1640.891927 {
				windD = [2][4]float64{
					{0.0, 10.0, 1389.84375, 1525.84375},
					{0.0, 10.0, 1640.891927, 1867.0},
				}
			} else if tomYPos <= 1958.915365 {
				windD = [2][4]float64{
					{0.0, 10.0, 1640.891927, 1867.0},
					{0.0, 10.0, 1958.915365, 2300.0},
				}
			} else {
				windD = [2][4]float64{
					{0.0, 10.0, 1958.915365, 2300.0},
					{0.0, 10.0, 2262.979167, 2710.0},
				}
			}
		} else {
			windOffset = math.Mod(wind, 15.0)
			windD = [2][4]float64{
				{0.0, 5.0, 1527.864583, 1471.875},
				{0.0, 5.0, 1791.894531, 1709.895833},
			}
		}

		windLow := windD[0][2] - ((windD[0][2]-windD[0][3])/(windD[0][1]-windD[0][0]))*math.Abs(windOffset)
		windHigh := windD[1][2] - ((windD[1][2]-windD[1][3])/(windD[1][1]-windD[1][0]))*math.Abs(windOffset)

		var windFactor float64
		if wind >= 0.0 {
			windFactor = (tomYPos - initialFactor[0][2]) / (initialFactor[1][2] - initialFactor[0][2])
		} else {
			windFactor = (tomYPos - windD[0][2]) / (windD[1][2] - windD[0][2])
		}

		windYPos = (windHigh-windLow)*windFactor + windLow
		windXPos = math.Abs(windXOffset)
	}

	return windXPos, windYPos
}

func (s CalculationsService) tickGenerator(min, max float64) plot.Ticker {
	return plot.ConstantTicks(func() []plot.Tick {
		var ticks []plot.Tick
		step := 20.0
		for val := min; val <= max; val += step {
			ticks = append(ticks, plot.Tick{
				Value: val,
				Label: fmt.Sprintf("%.0f", val),
			})
		}
		return ticks
	}())
}

func (s CalculationsService) LandingDistance(oat temperature.Temperature, pa pressure.Altitude, tow weight_balance.Mass, wi wind.Wind, chartType calculations.ChartType) (io.Reader, float64, float64, error) {
	oatXStart := 562.923177
	oatXEnd := 1870.93099
	oatXUnits := 70.0

	var w float64
	if wi.Direction == wind.DirectionTailwind {
		w = -wi.Speed.Float64()
	} else {
		w = wi.Speed.Float64()
	}

	oatY := [][2]interface{}{
		{0.0, []float64{1902.34, 1948.34, 1994.34, 2042.32, 2090.33, 2136.33, 2184.34, 2234.34}},
		{2000.0, []float64{2002.34, 2054.33, 2104.33, 2158.33, 2210.32, 2262.34, 2316.34, 2370.34}},
		{4000.0, []float64{2114.32, 2172.33, 2228.32, 2286.33, 2344.34, 2404.33, 2462.34, 2522.33}},
		{6000.0, []float64{2242.32, 2304.33, 2368.33, 2432.32, 2498.34, 2562.33, 2628.32, 2694.34}},
		{8000.0, []float64{2384.34, 2454.33, 2526.33, 2598.34, 2670.34, 2742.32, 2814.33, 2888.31}},
	}

	yBracket := int((oat.Float64() + 30) / 10)
	yInterpolated1, yInterpolated2, err := s.interpolateYValues(pa.Float64(), oatY, yBracket)

	if err != nil {
		return nil, 0, 0, err
	}

	pOatX := (oatXEnd - oatXStart) / oatXUnits
	yOffset := math.Mod(oat.Float64()+30.0, 10)
	pOatY := (yInterpolated2 - yInterpolated1) / 10

	tomXStart := 2077.115885
	tomXEnd := 3263.216146
	tomUnits := 750.0 - 550.0
	tomXOffset := (750.0-tow.Kilo())*(tomXEnd-tomXStart)/tomUnits + tomXStart

	tom := [2][4]float64{
		{0.0, 200.0, 1906.05, 1796.06},
		{0.0, 200.0, 2002.08, 1882.06},
	}

	tomYPos := s.interpolateTomY(tom, tow.Kilo(), yInterpolated1, pOatY, yOffset)

	windXStart := 3439.388021
	windXEnd := 4933.561198
	windXPos, windYPos := s.calculateWindPosition(windXStart, windXEnd, w, tomYPos, tomXOffset)

	obs := [2][2]float64{
		{1467.55, 1171.48},
		{1631.61, 1241.50},
	}

	perfYStart := 965.46224
	perfYEnd := 3261.946615
	perfUnits := 1000.0

	obsYPos := s.calculateObstacleY(windYPos, obs)

	ldrGRValues := []float64{
		1171.484375, 1241.503906, 1309.53776, 1379.557292, 1447.558594,
		1517.578125, 1585.579427, 1653.613281, 1723.632813, 1791.634115,
		1861.653646, 1929.654948, 1999.674479,
	}
	ldrDRValues := []float64{
		1467.545573, 1631.608073, 1797.65625, 1961.686198, 2125.716146,
		2289.746094, 2453.776042, 2617.80599, 2781.835938, 2947.884115,
		3111.914063, 3275.94401, 3440.00651,
	}

	ldrGR := s.findNextValue(ldrGRValues, obsYPos, perfYEnd)
	ldrDR := s.findNextValue(ldrDRValues, windYPos, perfYEnd)

	ldgGR := (ldrGR - perfYStart) / (perfYEnd - perfYStart) * perfUnits
	ldgDR := (ldrDR - perfYStart) / (perfYEnd - perfYStart) * perfUnits

	tmpl, err := template.New("svgTemplate").Parse(s.ldrChart.String())
	if err != nil {
		return nil, 0, 0, err
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, LdrData{
		OatX:  fmt.Sprintf("%.5f", oatXStart+pOatX*(oat.Float64()+30.0)),
		OatY:  fmt.Sprintf("%.5f", yInterpolated1+pOatY*yOffset),
		LmX:   fmt.Sprintf("%.5f", tomXOffset),
		LmY:   fmt.Sprintf("%.5f", tomYPos),
		WindX: fmt.Sprintf("%.5f", windXPos),
		WindY: fmt.Sprintf("%.5f", windYPos),
		ObY:   fmt.Sprintf("%.5f", obsYPos),
	})

	if err != nil {
		return nil, 0, 0, err
	}

	if chartType == calculations.ChartTypePNG {
		png, err := s.imageService.SvgToPng(&output)
		if err != nil {
			return nil, 0, 0, err
		}

		return png, ldgDR, ldgGR, nil
	}

	return &output, ldgDR, ldgGR, nil
}

func (s CalculationsService) WeightAndBalance(callSign callsign.CallSign, takeOffMassMoment weight_balance.MassMoment, landingMassMoment weight_balance.MassMoment, withinLimits bool, chartType calculations.ChartType) (io.Reader, error) {
	p := plot.New()

	p.Title.Text = fmt.Sprintf("Gewicht en Balans (%s)", callSign.String())
	p.X.Label.Text = "Mass Moment [kg m]"
	p.Y.Label.Text = "Mass [kg]"

	p.X.Min = 230
	p.X.Max = 430
	p.Y.Min = 550
	p.Y.Max = 770

	p.X.Tick.Marker = s.tickGenerator(p.X.Min, p.X.Max)
	p.Y.Tick.Marker = s.tickGenerator(p.Y.Min, p.Y.Max)

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
		return nil, err
	}

	polygon.LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	polygon.Color = color.RGBA{R: 255, G: 0, B: 0, A: 64}
	p.Add(polygon)

	points := plotter.XYs{
		{X: takeOffMassMoment.KGM(), Y: takeOffMassMoment.Mass.Kilo()},
	}

	points2 := plotter.XYs{
		{X: landingMassMoment.KGM(), Y: landingMassMoment.Mass.Kilo()},
	}

	scatterTakeOff, err := plotter.NewScatter(points)
	if err != nil {
		return nil, err
	}

	scatterTakeOff.GlyphStyle.Shape = draw.CircleGlyph{}
	if withinLimits {
		scatterTakeOff.GlyphStyle.Color = color.RGBA{G: 255, A: 255}
	} else {
		scatterTakeOff.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	}

	scatterTakeOff.GlyphStyle.Radius = vg.Points(5)

	p.Add(scatterTakeOff)

	scatterLanding, err := plotter.NewScatter(points2)
	if err != nil {
		return nil, err
	}

	scatterLanding.GlyphStyle.Shape = draw.CircleGlyph{}
	scatterLanding.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	scatterLanding.GlyphStyle.Radius = vg.Points(5)

	p.Add(scatterLanding)

	p.Legend.Top = false
	p.Legend.XOffs = vg.Centimeter * -1.5
	p.Legend.YOffs = vg.Centimeter * 1.5
	p.Legend.Padding = vg.Points(5)

	p.Legend.Add("CG Envelope", polygon)
	p.Legend.Add("Take-off Point", scatterTakeOff)
	p.Legend.Add("Landing Point", scatterLanding)

	var buf bytes.Buffer

	switch chartType {
	case calculations.ChartTypeSVG:

		c := canvas.New(150.0, 150.0)
		gCanvas := renderers.NewGonumPlot(c)

		p.Draw(gCanvas)

		err = c.Write(&buf, renderers.SVG())

		if err != nil {
			return nil, err
		}

		return &buf, nil
	default:
		writer, err := p.WriterTo(8*vg.Inch, 8*vg.Inch, "png")
		if err != nil {
			return nil, err
		}

		if _, err := writer.WriteTo(&buf); err != nil {
			return nil, err
		}

		return &buf, nil
	}
}

func (s CalculationsService) TakeOffDistance(oat temperature.Temperature, pa pressure.Altitude, tow weight_balance.Mass, wi wind.Wind, chartType calculations.ChartType) (io.Reader, float64, float64, error) {
	temp := oat.Float64()
	pressureAltitude := pa.Float64()
	mtow := tow.Kilo()
	w := wi.Speed.Float64()

	if wi.Direction == wind.DirectionTailwind {
		w = -w
	}

	oatXStart := 562.923177
	oatXEnd := 2168.91276
	oatXUnits := 70.0

	type OATData struct {
		PressureAltitude float64
		Values           [8]float64
	}

	oatY := []OATData{
		{0.0, [8]float64{1614.322917, 1656.315104, 1698.339844, 1742.317708, 1788.313802, 1834.342448, 1882.324219, 1932.324219}},
		{2000.0, [8]float64{1702.34375, 1750.325521, 1800.325521, 1850.325521, 1902.34375, 1956.315104, 2010.31901, 2066.341146}},
		{4000.0, [8]float64{1804.329427, 1860.31901, 1916.341146, 1974.316406, 2034.342448, 2096.321615, 2160.31901, 2224.316406}},
		{6000.0, [8]float64{1924.316406, 1988.313802, 2052.34375, 2120.345052, 2190.332031, 2262.33724, 2334.342448, 2410.31901}},
		{8000.0, [8]float64{2064.322917, 2138.313802, 2214.322917, 2292.317708, 2372.330729, 2456.315104, 2540.332031, 2628.320313}},
	}

	yBracket := int((oat + 30.0) / 10.0)

	var yInterpolated [2]float64
	var y0, y1 OATData

	if pressureAltitude <= 2000.0 {
		y0 = oatY[0]
		y1 = oatY[1]
	} else if pressureAltitude <= 4000.0 {
		y0 = oatY[1]
		y1 = oatY[2]
	} else if pressureAltitude <= 6000.0 {
		y0 = oatY[2]
		y1 = oatY[3]
	} else if pressureAltitude <= 8000.0 {
		y0 = oatY[3]
		y1 = oatY[4]
	} else {
		panic("Drukhoogte niet binnen bereik")
	}

	yFactor := (pressureAltitude - y0.PressureAltitude) / (y1.PressureAltitude - y0.PressureAltitude)
	yInterpolated[0] = s.interpolate(y0.Values[yBracket], y1.Values[yBracket], yFactor)
	yInterpolated[1] = s.interpolate(
		y0.Values[min(yBracket+1, len(y0.Values)-1)],
		y1.Values[min(yBracket+1, len(y1.Values)-1)],
		yFactor,
	)

	pOatX := (oatXEnd - oatXStart) / oatXUnits
	yOffset := math.Mod(temp+30.0, 10.0)
	pOatY := (yInterpolated[1] - yInterpolated[0]) / 10.0

	tomXStart := 2367.122396
	tomXEnd := 3777.246094
	tomUnits := 750.0 - 550.0
	tomXOffset := (750.0-mtow)*((tomXEnd-tomXStart)/tomUnits) + tomXStart

	tomData := [2][4]float64{
		{0.0, 200.0, 1632.03125, 1400.032552},
		{0.0, 200.0, 1718.033854, 1454.003906},
	}
	tomY := [2]float64{
		s.interpolate(tomData[0][2], tomData[0][3], (750.0-mtow)/(tomData[0][1]-tomData[0][0])),
		s.interpolate(tomData[1][2], tomData[1][3], (750.0-mtow)/(tomData[1][1]-tomData[1][0])),
	}
	tomYPos := s.interpolate(
		tomY[0],
		tomY[1],
		(yInterpolated[0]+(pOatY*yOffset)-tomData[0][2])/(tomData[1][2]-tomData[0][2]),
	)

	windXStart := 3965.429687
	windXEnd := 5211.621094
	windXPos, windYPos := s.calculateWindPositionTOD(windXStart, windXEnd, w, tomYPos, tomXOffset)

	obsData := [2][2]float64{
		{1395.703125, 1727.766927},
		{1491.731771, 1905.794271},
	}
	obsYPos := s.interpolateObstacleY(windYPos, obsData)

	perfYStart := 1009.635417
	perfYEnd := 4222.200521
	perfUnits := 1000.0

	torGRValues := []float64{
		1395.703125, 1491.731771, 1587.727865, 1683.75651, 1779.785156,
		1877.799479, 1973.795573, 2069.824219, 2165.852865, 2261.848958,
		2359.895833, 2455.891927, 2551.920573, 2655.924479,
	}
	torDRValues := []float64{
		1727.766927, 1905.794271, 2085.839844, 2265.852865, 2443.880208,
		2623.925781, 2803.938802, 2983.984375, 3162.011719, 3342.057292,
		3522.070313, 3700.097656, 3880.143229, 4076.171875,
	}

	torGR := s.findNextValue(torGRValues, windYPos, perfYEnd)
	torDR := s.findNextValue(torDRValues, obsYPos, perfYEnd)

	todGR := (torGR - perfYStart) / (perfYEnd - perfYStart) * perfUnits
	todDR := (torDR - perfYStart) / (perfYEnd - perfYStart) * perfUnits

	oatXBase := oatXStart + (pOatX * (temp + 30.0))
	oatYBase := yInterpolated[0] + (pOatY * yOffset)

	tmpl, err := template.New("svgTemplate").Parse(s.tdrChart.String())
	if err != nil {
		return nil, 0, 0, err
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, TodData{
		OatX:  fmt.Sprintf("%.5f", oatXBase),
		OatY:  fmt.Sprintf("%.5f", oatYBase),
		TomX:  fmt.Sprintf("%.5f", tomXOffset),
		TomY:  fmt.Sprintf("%.5f", tomYPos),
		WindX: fmt.Sprintf("%.5f", windXPos),
		WindY: fmt.Sprintf("%.5f", windYPos),
		ObsY:  fmt.Sprintf("%.5f", obsYPos),
	})
	if err != nil {
		return nil, 0, 0, err
	}

	if chartType == calculations.ChartTypePNG {
		png, err := s.imageService.SvgToPng(&output)
		if err != nil {
			return nil, 0, 0, err
		}

		return png, todGR, todDR, nil
	}

	return &output, todGR, todDR, nil
}

func (s CalculationsService) Calculations(callSign *callsign.CallSign, pilot *weight_balance.Mass, pilotSeat *seat.Position, passenger *weight_balance.Mass, passengerSeat *seat.Position, baggage *weight_balance.Mass, outsideAirTemperature *temperature.Temperature, pa *pressure.Altitude, w *wind.Wind, f *fuel.Fuel, tripDuration *time.Duration, alternateDuration *time.Duration) (*calculations.Calculations, error) {
	sheet := &calculations.Calculations{}

	sheet.CallSign = *callSign
	sheet.PressureAltitude = *pa
	sheet.OAT = *outsideAirTemperature
	sheet.Wind = *w

	var err error
	sheet.TakeOffWeightAndBalance, err = calculations.NewWeightAndBalance(*callSign, *pilot, *pilotSeat, passenger, passengerSeat, baggage, *f)
	if err != nil {
		return sheet, err
	}

	sheet.FuelPlanning, err = calculations.NewFuelPlanning(*tripDuration, *alternateDuration, *f, f.Volume.Type)
	if err != nil {
		return sheet, err
	}

	sheet.LandingWeightAndBalance, err = calculations.NewWeightAndBalance(*callSign, *pilot, *pilotSeat, passenger, passengerSeat, baggage, fuel.Subtract(*f, sheet.FuelPlanning.Trip, sheet.FuelPlanning.Taxi))
	if err != nil {
		return sheet, err
	}

	_, todRR, todDR, err := s.TakeOffDistance(*outsideAirTemperature, *pa, sheet.TakeOffWeightAndBalance.Total.Mass, *w, calculations.ChartTypeSVG)

	if err != nil {
		return sheet, err
	}

	_, ldrDR, ldrGR, err := s.LandingDistance(*outsideAirTemperature, *pa, sheet.LandingWeightAndBalance.Total.Mass, *w, calculations.ChartTypeSVG)
	if err != nil {
		return sheet, err
	}
	sheet.Performance = &calculations.Performance{
		TakeOffRunRequired:        todRR,
		TakeOffDistanceRequired:   todDR,
		LandingDistanceRequired:   ldrDR,
		LandingGroundRollRequired: ldrGR,
	}

	return sheet, nil
}
