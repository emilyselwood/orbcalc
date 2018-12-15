// A horibly simple web server that will provide a set of point coords for all the asteroids

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wselwood/gompcreader"
	"github.com/wselwood/orbcalc/orbconvert"
	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
)

const dataStorePath = "static/data/"
const dataFile = dataStorePath + "data-%v.csv"
const numFiles = 16

var generate = flag.Bool("gen", false, "Should the data file be generated")
var dataPath = flag.String("data", "", "path to the minor planet center data file")

func main() {
	flag.Parse()

	if *generate && *dataPath == "" {
		flag.Usage()
		log.Fatal("http parameter is required")
	}

	if *generate {
		generateData()
		generateMajorPlanetData()
	}

	log.Println("starting server")
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Fatal(http.ListenAndServe(":8000", nil))
}

// internal memory cache for the positions
var positions []*orbcore.Position

func generateData() {

	log.Println("generating data...")

	mpcReader, err := gompcreader.NewMpcReader(*dataPath)
	if err != nil {
		log.Fatal("error creating mpcReader ", err)
	}
	defer mpcReader.Close()
	outFiles := make([]*os.File, numFiles)
	for i := 0; i < numFiles; i++ {
		out, err := os.Create(filename(i))
		if err != nil {
			log.Fatal("error opening output file", err)
		}
		outFiles[i] = out
	}

	defer func() {
		for _, f := range outFiles {
			f.Close()
		}
	}()

	var count int
	result, err := mpcReader.ReadEntry()
	for err == nil {

		orb := orbconvert.ConvertFromMinorPlanet(result)
		pos := orbcore.OrbitToPosition(orb)

		fmt.Fprintf(outFiles[count%numFiles], "%v,%v,%v\n", pos.X, pos.Y, pos.Z)
		count++
		result, err = mpcReader.ReadEntry()
	}

	if err != nil && err != io.EOF {
		log.Fatal("error reading", err)
	}

	log.Println("Loaded", count, "records")
}

func filename(i int) string {
	return fmt.Sprintf(dataFile, i)
}

func generateMajorPlanetData() {
	for _, p := range orbdata.SolarSystem {

		if err := fileForPlanet(p); err != nil {
			log.Fatal(err)
		}
	}
}

func fileForPlanet(orb orbcore.Orbit) error {

	out, err := os.Create(dataStorePath + orb.ID + ".csv")
	if err != nil {
		log.Fatal("error opening output file", err)
	}
	defer out.Close()

	positions := orbcore.MeanMotionFullOrbit(&orb, 366)
	for _, pos := range positions {
		p := orbcore.OrbitToPosition(pos)
		fmt.Fprintf(out, "%v,%v,%v\n", p.X, p.Y, p.Z)
	}

	return nil
}
