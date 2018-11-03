package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"time"
	"io"

	"github.com/wselwood/gompcreader"
	"github.com/wselwood/orbcalc/orbconvert"
	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
	"gonum.org/v1/gonum/mat"

	"github.com/paulbellamy/ratecounter"
)

const monitoringInterval = 1 * time.Second
const processors = 1
const channelSize = 10000

var inputfile = flag.String("in", "", "the minor planet center file to read")
var outputfile = flag.String("out", "", "path to output file")
var count = flag.Int("count", 1000000, "number of records to run")
var skip = flag.Int("skip", 0, "number of records from the begining to skip")
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

	// rate counters for each stage
	counter1 := ratecounter.NewRateCounter(monitoringInterval)
	counter2 := ratecounter.NewRateCounter(monitoringInterval)
	counter3 := ratecounter.NewRateCounter(monitoringInterval)

	timer := time.NewTicker(monitoringInterval)
	defer timer.Stop()

	go func() {
		for range timer.C {
			log.Printf("stage1: %v stage2: %v, stage3: %v", counter1.Rate(), counter2.Rate(), counter3.Rate())
		}
	}()

	stage1 := make(chan *orbcore.Orbit, channelSize)
	stage2 := make(chan *orbcore.Orbit, channelSize)

	var readGroup sync.WaitGroup
	var fanGroup sync.WaitGroup
	var complete sync.WaitGroup

	readGroup.Add(1)
	go stageRead(*inputfile, *count, *skip, stage1, &readGroup, counter1)

	for i := 0; i < processors; i++ {
		fanGroup.Add(1)
		go stageMeanMotion(stage1, stage2, &fanGroup, counter2)
	}

	complete.Add(1)
	go stageOutput(*outputfile, stage2, &complete, counter3)

	readGroup.Wait()
	log.Println("done waiting for read")

	fanGroup.Wait()
	close(stage2)

	log.Println("done waiting for fan")

	complete.Wait()
	log.Println("done")
	
}

func stageRead(inputfile string, target int, skip int, output chan *orbcore.Orbit, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	mpcReader, err := gompcreader.NewMpcReader(inputfile)
	if err != nil {
		log.Fatal("error creating mpcReader ", err)
	}
	defer mpcReader.Close()
	defer close(output)
	defer wg.Done()

	var count int
	result, err := mpcReader.ReadEntry()
	for err == nil {
		orb := orbconvert.ConvertFromMinorPlanet(result)
		if skip == 0 {
			//fmt.Println(orb)
			output <- orb
		} else {
			skip--
		}
		counter.Incr(1)
		count++
		if count >= target {
			return
		}
		result, err = mpcReader.ReadEntry()
	}
	if err != nil && err != io.EOF {
		log.Fatal("error reading", err)
	}

}

func stageMeanMotion(in chan *orbcore.Orbit, output chan *orbcore.Orbit, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()
	
	for orb := range in {
		output <- orbcore.MeanMotion(orbdata.SunGrav, orb, 1*24*60*60)
		counter.Incr(1)
	}
	
}

func stageOutput(outputPath string, in chan *orbcore.Orbit, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()
	
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("error creating outputfile ", err)
	}
	defer f.Close()
	for orb := range in {
		r, _ := orbcore.OrbitToVector(orb)
		fmt.Fprintf(f, "%s,%v\n", orb.ID, mat.Formatted(r.T(), mat.Prefix(""), mat.Squeeze()))
		counter.Incr(1)
	}

	
}
