package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/emilyselwood/gompcreader"
	"github.com/emilyselwood/orbcalc/orbconvert"
	"github.com/emilyselwood/orbcalc/orbcore"

	"github.com/paulbellamy/ratecounter"
)

const monitoringInterval = 1 * time.Second
const processors = 3
const channelSize = 100000

var inputfile = flag.String("in", "", "the minor planet center file to read")
var outputfile = flag.String("out", "", "path to output file")
var count = flag.Int("count", 1000000, "number of records to run")
var skip = flag.Int("skip", 0, "number of records from the begining to skip")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

/*
An example program that uses the gompcreader and calculates the position in space for each object.
This example is a little more complex than strictly needed as it parallises the processing and keeps counters.
*/
func main() {
	flag.Parse()

	// If a cpu project has been requested then set that up.
	// This is not strictly needed but it does help us know what is going on with the program
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

	// rate counters for each processing stage
	counter1 := ratecounter.NewRateCounter(monitoringInterval)
	counter2 := ratecounter.NewRateCounter(monitoringInterval)
	counter3 := ratecounter.NewRateCounter(monitoringInterval)

	// start an thing to print rate information as we are processing.
	timer := time.NewTicker(monitoringInterval)
	defer timer.Stop()

	go func() {
		for range timer.C {
			log.Printf("stage1: %v stage2: %v stage3: %v", counter1.Rate(), counter2.Rate(), counter3.Rate())
		}
	}()

	// Prepares channels to transfer data between stages of the processing.
	stage1 := make(chan *orbcore.Orbit, channelSize)
	stage2 := make(chan *orbcore.Position, channelSize)

	// Wait groups for each stage so we know when things are done.
	var readGroup sync.WaitGroup
	var fanGroup sync.WaitGroup
	var complete sync.WaitGroup

	// Setup stage one which reads in the input file, parses the records from it and passes them on to the next stage.
	readGroup.Add(1)
	go stageRead(*inputfile, *count, *skip, stage1, &readGroup, counter1)

	// Stage two progates an object forward one day and then passes it on.
	for i := 0; i < processors; i++ {
		fanGroup.Add(1)
		go stageMeanMotion(stage1, stage2, &fanGroup, counter2)
	}

	// The final stage converts the orbit information into a position in space and then writes it to a file.
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

// stageRead opens a file using the gompcreader project and reads out orbital information.
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
	if err != io.EOF {
		log.Fatal("error reading", err)
	}

}

func stageMeanMotion(in chan *orbcore.Orbit, output chan *orbcore.Position, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()
	oneDay := 24 * time.Hour
	for orb := range in {
		log.Println(orb.Epoch)

		r := orbcore.MeanMotionStepped(orb, oneDay, 2000)
		for _, o := range r {
			output <- orbcore.OrbitToPosition(o)
		}
		counter.Incr(1)
	}
}

func stageOutput(outputPath string, in chan *orbcore.Position, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("error creating outputfile ", err)
	}
	defer f.Close()

	w := bufio.NewWriterSize(f, 64*1024)
	defer w.Flush()
	for orb := range in {
		w.WriteString(orb.String())
		w.WriteRune('\n')
		counter.Incr(1)
	}

}
