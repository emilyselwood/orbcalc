package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type row struct {
	ID  string
	Day int
	X   float64
	Y   float64
	Z   float64
}

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

	file, err := os.Open(*in)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	count := 0
	result := make([]*row, 366) // there should be this many entries in the file.
	for scanner.Scan() {
		r, err := parseLine(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		result[count] = r
		count++
	}

	if err := scanner.Err(); err != nil {
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

	if err := p.Save(4*vg.Inch, 4*vg.Inch, *out); err != nil {
		log.Fatal(err)
	}
}

func rowsToPointsXY(rows []*row) plotter.XYs {
	pts := make(plotter.XYs, len(rows))
	for i := range pts {
		pts[i].X = rows[i].X
		pts[i].Y = rows[i].Y
	}
	return pts
}

func parseLine(line string) (*row, error) {
	parts := strings.Split(line, ",")

	i, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	x, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, err
	}
	y, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, err
	}
	z, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, err
	}

	result := row{
		ID:  parts[0],
		Day: i,
		X:   x,
		Y:   y,
		Z:   z,
	}
	return &result, nil
}
