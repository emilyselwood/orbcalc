package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	out := flag.String("out", "out.png", "output filename for plot")

	flag.Parse()

	if *out == "" {
		log.Fatal("Need an output filename")
	}

	// now to plot
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Title.Text = fmt.Sprintf("Solar system")

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	for i, orb := range orbdata.SolarSystem {
		process(p, orb, rainbow(len(orbdata.SolarSystem), i))
	}

	if err := p.Save(800, 800, *out); err != nil {
		log.Fatal(err)
	}
}

func process(p *plot.Plot, orb orbcore.Orbit, c color.RGBA) {

	result := propogate(&orb, orbdata.SunGrav)

	l, err := plotter.NewLine(rowsToPointsXY(result))
	if err != nil {
		log.Fatal(err)
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = c

	p.Add(l)

}

func rowsToPointsXY(rows []*orbcore.Position) plotter.XYs {
	pts := make(plotter.XYs, len(rows))
	for i := range pts {
		pts[i].X = rows[i].X
		pts[i].Y = rows[i].Y
	}
	return pts
}

func propogate(orb *orbcore.Orbit, parent float64) []*orbcore.Position {
	steps := orbcore.MeanMotionFullOrbit(parent, orb, 366)
	result := make([]*orbcore.Position, len(steps))
	for i, d := range steps {
		result[i] = orbcore.OrbitToPosition(d, parent)
	}
	return result
}

type shortTicks struct{}

// Makes the ticks shortened.
func (shortTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		if strings.HasSuffix(t.Label, "e+08") {
			tks[i].Label = t.Label[:len(t.Label)-4]
		}
	}
	return tks
}

func rainbow(numOfSteps, step int) color.RGBA {

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
