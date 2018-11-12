package main

import (
	"flag"
	"log"

	"github.com/wselwood/orbcalc/orbcore"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
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

	p.Title.Text = result[0].ID
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	if err := plotutil.AddLinePoints(p, rowsToPointsXY(result)); err != nil {
		log.Fatal(err)
	}

	if err := p.Save(8*vg.Inch, 8*vg.Inch, *out); err != nil {
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
