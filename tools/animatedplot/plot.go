package main

import (
	
	"flag"
	"io"
	"log"
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"time"
	"strconv"
	"path/filepath"

	"github.com/wselwood/gompcreader"
	"github.com/wselwood/orbcalc/orbconvert"
	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"

	"github.com/paulbellamy/ratecounter"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	
)

const monitoringInterval = 1 * time.Second
const processors = 3
const channelSize = 100000

var inputfile = flag.String("in", "", "the minor planet center file to read")
var outputPath = flag.String("out", "", "path to output files")
var count = flag.Int("count", 6000, "number of frames to run")
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

	if *outputPath == "" {
		log.Fatal("No output file prvided. Use the -out /path/to/outputfiles")
	}

	os.MkdirAll(*outputPath, os.ModePerm)

	saveChan := make(chan *savePack, 100)
	var saveWait sync.WaitGroup
	for i := 0; i < 4; i++ {
		saveWait.Add(1)
		go func() {
			defer saveWait.Done()
			for s := range saveChan {		
				if err := s.plot.Save(1000, 1000, s.filename); err != nil {
					panic(err)
				}
			}
		}()
	}

	for i := int64(0); i < int64(*count); i++ {
		processFrame(i, saveChan)
	}

	close(saveChan)
	saveWait.Wait()
}

func processFrame(days int64, saveChan chan *savePack) {
	// rate counters for each processing stage
	counter1 := ratecounter.NewRateCounter(monitoringInterval)
	counter2 := ratecounter.NewRateCounter(monitoringInterval)
	counter3 := ratecounter.NewRateCounter(monitoringInterval)
	counter4 := ratecounter.NewRateCounter(monitoringInterval)

	// start an thing to print rate information as we are processing.
	timer := time.NewTicker(monitoringInterval)
	defer timer.Stop()

	go func() {
		for range timer.C {
			log.Printf("stage1: %v stage2: %v stage3: %v stage4: %v", counter1.Rate(), counter2.Rate(), counter3.Rate(), counter4.Rate())
		}
	}()

	// Prepares channels to transfer data between stages of the processing.
	stage1 := make(chan *orbcore.Orbit, channelSize)
	stage2 := make(chan *orbcore.Orbit, channelSize)
	stage3 := make(chan *orbcore.Position, channelSize)

	// Wait groups for each stage so we know when things are done.
	var readGroup sync.WaitGroup
	var meanMotionGroup sync.WaitGroup
	var positionGroup sync.WaitGroup
	var complete sync.WaitGroup

	// Setup stage one which reads in the input file, parses the records from it and passes them on to the next stage.
	readGroup.Add(1)
	go stageRead(*inputfile, 1000000, 0, stage1, &readGroup, counter1)

	// Stage two progates an object forward one day and then passes it on.
	for i := 0; i < processors; i++ {
		meanMotionGroup.Add(1)
		go stageMeanMotion(days, stage1, stage2, &meanMotionGroup, counter2)

		positionGroup.Add(1)
		go stagePosition(stage2, stage3, &positionGroup, counter3)
	}

	// The final stage converts the orbit information into a position in space and then writes it to a file.
	complete.Add(1)
	go stagePlot(days, stage3, saveChan, *outputPath, &complete, counter4)

	readGroup.Wait()

	meanMotionGroup.Wait()
	close(stage2)

	positionGroup.Wait()
	close(stage3)

	complete.Wait()
	log.Println("done", days)
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
	if err != nil && err != io.EOF {
		log.Fatal("error reading", err)
	}

}

func stageMeanMotion(days int64, in chan *orbcore.Orbit, output chan *orbcore.Orbit, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()
	offset := 24 * time.Hour *time.Duration(days)
	for orb := range in {
		output <- orbcore.MeanMotion(orbdata.SunGrav, orb, offset) 
		counter.Incr(1)
	}
}

func stagePosition(in chan *orbcore.Orbit, out chan *orbcore.Position, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()
	for orb := range in {
		out <- orbcore.OrbitToPosition(orb)
		counter.Incr(1)
	}
}

func stagePlot(days int64, in chan *orbcore.Position, saveChan chan *savePack, outputPath string, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()

	
	var data plotter.XYs
	for orb := range in {
		data = append(data, struct{X, Y float64}{orb.X, orb.Y})
		counter.Incr(1)
	}

	p := createPlot(days)
	scatter, err := plotter.NewScatter(data)
	if err != nil {
		panic(err)
	}
	scatter.Radius = vg.Points(1)
	
	p.Add(scatter)

	filename := filepath.Join(outputPath, fmt.Sprintf("frame_%05d.png", days))
	saveChan <- &savePack{
		filename:filename,
		plot:p,
	}
	
}

func createPlot(t int64) *plot.Plot {
	
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = strconv.FormatInt(t, 10)
	p.X.Label.Text = "X (km)"
	p.Y.Label.Text = "Y (km)"

	return p
}

type savePack struct {
	filename string
	plot *plot.Plot
}