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
)

var inputfile = flag.String("in", "", "the minor planet center file to read")
var untilId = flag.String("until", "100000", "id to stop when found")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

/*
An example program that uses the gompcreader and calculates the position in space for each object.
*/
func main() {
	// TODO: cpu profile
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
	var total float64
	var count int64
	result, err := mpcReader.ReadEntry()
	for err == nil {
		orb := orbconvert.ConvertFromMinorPlanet(result)
		r, v := orbcore.OrbitToVector(orb)
		if *cpuprofile != "" {
			total = total + r.At(2, 0) + v.At(2, 0)
		} else {
			fmt.Printf("%s:%s\n", result.ID, result.ReadableDesignation)
			fmt.Println("r:", r)
			fmt.Println("v:", v)
		}

		count = count + 1
		if result.ID == *untilId {
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
