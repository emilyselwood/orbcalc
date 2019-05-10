package orbplot

import (
	"image/color"
	"math"

	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// PlotSolarSystemLines plots the major planets of the solar system on the provided plot
func PlotSolarSystemLines(p *plot.Plot, legend bool) error {
	if err := PlotFullOrbitLines(p, orbdata.SolarSystem, RainbowList(len(orbdata.SolarSystem)), legend); err != nil {
		return err
	}
	return PlotSun(p)
}

// PlotInnerSolarSystemLines plots the major planets of the inner solar system on the provided plot
func PlotInnerSolarSystemLines(p *plot.Plot, legend bool) error {
	if err := PlotFullOrbitLines(p, orbdata.InnerSolarSystem, RainbowList(len(orbdata.InnerSolarSystem)), legend); err != nil {
		return err
	}
	return PlotSun(p)
}

// PlotOuterSolarSystemLines plots the major planets of the outer solar system on the provided plot
func PlotOuterSolarSystemLines(p *plot.Plot, legend bool) error {
	if err := PlotFullOrbitLines(p, orbdata.OuterSolarSystem, RainbowList(len(orbdata.OuterSolarSystem)), legend); err != nil {
		return err
	}
	return PlotSun(p)
}

// PlotFullOrbitLines plots the full orbits of the provided orbits on the plot
func PlotFullOrbitLines(p *plot.Plot, orbits []orbcore.Orbit, colors []color.RGBA, legend bool) error {
	for i, orb := range orbits {
		if err := PlotFullOrbitLine(p, orb, colors[i], legend); err != nil {
			return err
		}
	}
	return nil
}

// PlotFullOrbitLine takes a plot and orbit and draws a line for its full orbit in the provided color
func PlotFullOrbitLine(p *plot.Plot, orb orbcore.Orbit, c color.RGBA, legend bool) error {

	result := propogate(&orb)

	l, err := plotter.NewLine(PositionToPointsXY(result))
	if err != nil {
		return err
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = c

	p.Add(l)
	if legend {
		p.Legend.Add(orb.ID, l)
	}
	return nil
}

// PlotPoints renders a point for each orbit object paired with the color list
func PlotPoints(p *plot.Plot, orb []*orbcore.Orbit, c []color.RGBA, legend bool) error {
	result := make([]*orbcore.Position, len(orb))
	for i, d := range orb {
		result[i] = orbcore.OrbitToPosition(d)
	}

	scatter, err := plotter.NewScatter(PositionToPointsXY(result))
	if err != nil {
		return err
	}

	scatter.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		var result draw.GlyphStyle

		result.Color = c[i]
		return result
	}

	p.Add(scatter)

	return nil
}

// PositionToPointsXY Converts many Position objects to something that can be plotted.
func PositionToPointsXY(rows []*orbcore.Position) plotter.XYs {
	pts := make(plotter.XYs, len(rows))
	for i := range pts {
		pts[i].X = rows[i].X
		pts[i].Y = rows[i].Y
	}
	return pts
}

// PlotSun adds a yellow dot at the origin.
func PlotSun(p *plot.Plot) error {
	pts := make(plotter.XYs, 1)
	pts[0].X = 0
	pts[0].Y = 0

	points, err := plotter.NewScatter(pts)
	if err != nil {
		return err
	}

	points.Color = color.RGBA{
		R: 255,
		G: 255,
		B: 0,
		A: 255,
	}

	points.Radius = 2
	points.Shape = draw.CircleGlyph{}

	p.Add(points)
	return nil
}

func propogate(orb *orbcore.Orbit) []*orbcore.Position {
	steps := orbcore.MeanMotionFullOrbit(orb, 366)
	result := make([]*orbcore.Position, len(steps))
	for i, d := range steps {
		result[i] = orbcore.OrbitToPosition(d)
	}
	return result
}

// RainbowList returns a list of colours
func RainbowList(numOfSteps int) []color.RGBA {
	result := make([]color.RGBA, numOfSteps)
	for i := 0; i < numOfSteps; i++ {
		result[i] = Rainbow(numOfSteps, i)
	}
	return result
}

// Rainbow generates a colour range that roughly matches the rainbow colour spectrum.
func Rainbow(numOfSteps, step int) color.RGBA {

	var r, g, b float64

	h := float64(step) / float64(numOfSteps)
	i := math.Floor(h * 6)
	f := h*6 - i
	q := 1 - f

	os := math.Remainder(i, 6)

	switch os {
	case 0:
		r = 255
		g = f * 255
		b = 0
	case 1:
		r = q * 255
		g = 255
		b = 0
	case 2:
		r = 0
		g = 255
		b = f * 255
	case 3:
		r = 0
		g = q * 255
		b = 255
	case 4:
		r = f * 255
		g = 0
		b = 255
	case 5:
		r = 255
		g = 0
		b = q * 255
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}
