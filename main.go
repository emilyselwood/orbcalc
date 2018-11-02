package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/wselwood/gompcreader"
	"github.com/wselwood/orbcalc/orbconvert"
	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
	"gonum.org/v1/gonum/mat"
)

var inputfile = flag.String("in", "", "the minor planet center file to read")
var outputfile = flag.String("out", "", "path to output file")
var count = flag.Int("count", 100000, "number of records to run")
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

	if *outputfile == "" {
		log.Fatal("No output file prvided. Use the -out /path/to/outputfile")
	}

	stage1 := make(chan *orbcore.Orbit, 1000)
	stage2 := make(chan *orbcore.Orbit, 1000)

	var fanGroup sync.WaitGroup
	var complete sync.WaitGroup

	go stageRead(*inputfile, *count, stage1)
	for i := 0; i < 4; i++ {
		fanGroup.Add(1)
		go stageMeanMotion(stage1, stage2, &fanGroup)
	}
	complete.Add(1)
	go stageOutput(*outputfile, stage2, &complete)

	fanGroup.Wait()
	close(stage2)

	complete.Wait()

}

func stageRead(inputfile string, target int, output chan *orbcore.Orbit) {
	mpcReader, err := gompcreader.NewMpcReader(inputfile)
	if err != nil {
		log.Fatal("error creating mpcReader ", err)
	}
	defer mpcReader.Close()
	var count int
	result, err := mpcReader.ReadEntry()
	for err == nil {
		orb := orbconvert.ConvertFromMinorPlanet(result)
		output <- orb
		count++
		if count >= target {
			break
		}
		result, err = mpcReader.ReadEntry()

	}
	close(output)
}

func stageMeanMotion(in chan *orbcore.Orbit, output chan *orbcore.Orbit, wg *sync.WaitGroup) {
	for orb := range in {
		output <- orbcore.MeanMotion(orbdata.SunGrav, orb, 1*(60*60*24))
	}
	wg.Done()
}

func stageOutput(outputPath string, in chan *orbcore.Orbit, wg *sync.WaitGroup) {
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("error creating outputfile ", err)
	}
	defer f.Close()
	for orb := range in {
		r, _ := orbcore.OrbitToVector(orb)
		fmt.Fprintf(f, "%s,%v\n", orb.ID, mat.Formatted(r.T(), mat.Prefix(""), mat.Squeeze()))
	}

	wg.Done()
}
