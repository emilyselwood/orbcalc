package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"github.com/wselwood/gompcreader"
	"github.com/wselwood/orbcalc/orbconvert"
	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
)

var inputfile = flag.String("in", "", "the minor planet center file to read")
var untilID = flag.String("until", "100000", "id to stop when found")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

/*
An example program that uses the gompcreader and calculates the position in space for each object.
*/
func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *inputfile == "" {
		log.Fatal("No input file provided. Use the -in /path/to/file")
	}

	mpcReader, err := gompcreader.NewMpcReader(*inputfile)
	if err != nil {
		log.Fatal("error creating mpcReader ", err)
	}
	defer mpcReader.Close()

	// Counters to make sure that the results don't get optomised away.
	// Paranoia is a hell of a thing
	var total float64
	var count int64

	// TODO: make this pipeline and threaded.
	// MeanMotion is slow as all things.

	result, err := mpcReader.ReadEntry()
	for err == nil {
		orb := orbconvert.ConvertFromMinorPlanet(result)
		orb = orbcore.MeanMotion(orbdata.SunGrav, orb, 1*(60*60*24))
		r, v := orbcore.OrbitToVector(orb)
		if *cpuprofile != "" {
			total = total + r.At(2, 0) + v.At(2, 0)
		} else {
			fmt.Printf("%s:%s\n", result.ID, result.ReadableDesignation)
			fmt.Println(orb.ID, "r:", r)
			fmt.Println(orb.ID, "v:", v)
		}

		count = count + 1
		if result.ID == *untilID {
			break
		}
		result, err = mpcReader.ReadEntry()
	}

	if err != nil && err != io.EOF {
		log.Fatal("error reading line\n", err)
	}

	log.Println("total:", total)
	log.Println("count:", count)
}

func stageRead() {

}

func stageMeanMotion() {

}

func stageOutput() {
	
}
