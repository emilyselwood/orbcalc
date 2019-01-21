// A horibly simple web server that will provide a set of point coords for all the asteroids

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	fs := http.FileServer(http.Dir("static"))
	
	http.HandleFunc("/", logWrapper(fs.ServeHTTP))

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func logWrapper(child func(rw http.ResponseWriter, req *http.Request)) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		child(rw, req)
		log.Println(req.Method, req.URL.Path, time.Since(startTime))
	}
}

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
			if err := f.Close(); err != nil {
				fmt.Println("Could not close file", err)
			}
		}
	}()

	var count int
	result, err := mpcReader.ReadEntry()
	for err == nil {

		orb := orbconvert.ConvertFromMinorPlanet(result)
		pos := orbcore.OrbitToPosition(orb)
		i := count % numFiles
		if _, err := fmt.Fprintf(outFiles[i], "%s,%v,%v,%v\n", pos.ID, pos.X, pos.Y, pos.Z); err != nil {
			log.Fatal("Could not write", err)
		}

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
	defer func(c *os.File) {
		if err := c.Close(); err != nil {
			fmt.Println("could not close", c, err)
		}
	}(out)

	positions := orbcore.MeanMotionFullOrbit(&orb, 366)
	for _, pos := range positions {
		p := orbcore.OrbitToPosition(pos)
		if _, err := fmt.Fprintf(out, "%v,%v,%v\n", p.X, p.Y, p.Z); err != nil {
			fmt.Println("Could not write file", err)
			return err
		}
	}

	return nil
}


type objectData struct {
	Epoch                       time.Time
	MeanAnomalyEpoch            float64 // nu
	ArgumentOfPerihelion        float64 // w argp
	LongitudeOfTheAscendingNode float64 // omega raan
	InclinationToTheEcliptic    float64 // i inc
	OrbitalEccentricity         float64 // e ecc
	MeanDailyMotion             float64
	SemimajorAxis               float64 // a p
	Orbit                       []point
}

type point struct {
	X float64
	Y float64
	Z float64
}
