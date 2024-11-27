package queries

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"io"
	"ppl-calculations/domain/calculations"
	"ppl-calculations/domain/callsign"
	"ppl-calculations/domain/weight_balance"
)

type WBChartHandler struct {
}

func NewWBChartHandler() WBChartHandler {
	return WBChartHandler{}
}

type WBChartRequest struct {
	CallSign          callsign.CallSign
	TakeOffMassMoment weight_balance.MassMoment
	LandingMassMoment weight_balance.MassMoment
	WithinLimits      bool
}

func (h WBChartHandler) tickGenerator(min, max float64) plot.Ticker {
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

func (h WBChartHandler) Handle(_ context.Context, request WBChartRequest) (io.Reader, error) {
	p := plot.New()

	p.Title.Text = fmt.Sprintf("Gewicht en Balans (%s)", request.CallSign.String())
	p.X.Label.Text = "Mass Moment [kg m]"
	p.Y.Label.Text = "Mass [kg]"

	p.X.Min = 230
	p.X.Max = 430
	p.Y.Min = 550
	p.Y.Max = 770

	p.X.Tick.Marker = h.tickGenerator(p.X.Min, p.X.Max)
	p.Y.Tick.Marker = h.tickGenerator(p.Y.Min, p.Y.Max)

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

	polygon.LineStyle.Color = color.RGBA{R: 255, A: 0}
	polygon.Color = color.RGBA{R: 255, A: 128}
	p.Add(polygon)

	points := plotter.XYs{
		{X: request.TakeOffMassMoment.KGM(), Y: request.TakeOffMassMoment.Mass.Kilo()},
	}

	points2 := plotter.XYs{
		{X: request.LandingMassMoment.KGM(), Y: request.LandingMassMoment.Mass.Kilo()},
	}

	scatterTakeOff, err := plotter.NewScatter(points)
	if err != nil {
		return nil, err
	}

	scatterTakeOff.GlyphStyle.Shape = draw.CircleGlyph{}
	if request.WithinLimits {
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

	c := canvas.New(150.0, 150.0)
	gCanvas := renderers.NewGonumPlot(c)

	p.Draw(gCanvas)

	var buf bytes.Buffer
	err = c.Write(&buf, renderers.SVG())
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
