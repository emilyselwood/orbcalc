package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/emilyselwood/orbcalc/orbplot"
	"gonum.org/v1/plot"
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

	p.X.Label.Text = "X (km)"
	p.Y.Label.Text = "Y (km)"

	if err := orbplot.PlotSolarSystemLines(p, true); err != nil {
		log.Fatal(err)
	}

	if err := p.Save(800, 800, *out); err != nil {
		log.Fatal(err)
	}
}
