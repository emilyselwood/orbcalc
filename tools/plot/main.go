package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/wselwood/orbcalc/orbcore"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	in := flag.String("in", "", "input csv file to plot")
	out := flag.String("out", "out.png", "output filename for plot")

	flag.Parse()

	if *in == "" {
		log.Fatal("Need an input filename")
	}

	if *out == "" {
		log.Fatal("Need an output filename")
	}

	result, err := orbcore.ReadPositionFile(*in)
	if err != nil {
		log.Fatal(err)
	}

	// now to plot
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = fmt.Sprintf("Orbit of %s between %v and %v", result[0].ID, result[0].Epoch, result[len(result)-1].Epoch)

	p.X.Label.Text = "X (1e8km)"
	p.X.Tick.Marker = shortTicks{}
	p.Y.Label.Text = "Y (1e8km)"
	p.Y.Tick.Marker = shortTicks{}

	l, err := plotter.NewLine(rowsToPointsXY(result))
	if err != nil {
		log.Fatal(err)
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l)

	if err := p.Save(800, 800, *out); err != nil {
		log.Fatal(err)
	}
}

func rowsToPointsXY(rows []*orbcore.Position) plotter.XYs {
	pts := make(plotter.XYs, len(rows))
	for i := range pts {
		pts[i].X = rows[i].X
		pts[i].Y = rows[i].Y
	}
	return pts
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
