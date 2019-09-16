// A horribly simple web server that will provide a set of point coords for all the asteroids

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/paulbellamy/ratecounter"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/emilyselwood/gompcreader"
	"github.com/emilyselwood/orbcalc/orbconvert"
	"github.com/emilyselwood/orbcalc/orbcore"
	"github.com/emilyselwood/orbcalc/orbdata"
)

const dataStorePath = "static/data/"
const dataFile = dataStorePath + "data-%v.csv"
const numFiles = 16
const monitoringInterval = 1 * time.Second
const processors = 3
const channelSize = 100000

var generate = flag.Bool("gen", false, "Should the data file be generated")
var dataPath = flag.String("data", "", "path to the minor planet center data file")
var genDate = flag.String("date", "2019-01-01", "Date to generate data for. YYYY-MM-DD format")

var asteroidData map[string]objectData

func main() {
	flag.Parse()

	if *generate && *dataPath == "" {
		flag.Usage()
		log.Fatal("data path parameter is required")
	}

	if *generate && *genDate == "" {
		flag.Usage()
		log.Fatal("generation date parameter is required")
	}

	if *generate {
		targetTime, err := time.Parse("2006-01-02", *genDate)
		if err != nil {
			flag.Usage()
			log.Fatal("generation date parameter is required", err)
		}

		generateData(targetTime)
		generateMajorPlanetData()
	}

	log.Println("loading data")
	if err := prepObjectData(*dataPath); err != nil {
		log.Fatal("Could not load data", err)
	}

	log.Println("starting server")
	fs := http.FileServer(http.Dir("static"))

	http.HandleFunc("/obj/", logWrapper(lookupAsteroid))
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

func lookupAsteroid(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.Path)
	id = strings.Replace(id, "+", " ", -1)
	v, ok := asteroidData[id]
	if !ok {
		rw.WriteHeader(404)
		return
	}
	if v.Orbit == nil {
		v.Orbit = make([]point, 366)
		for i, p := range orbcore.MeanMotionFullOrbit(v.toOrbit(), 365) {
			pos := orbcore.OrbitToPosition(p)
			v.Orbit[i].X = pos.X
			v.Orbit[i].Y = pos.Y
			v.Orbit[i].Z = pos.Z
		}
	}

	if err := json.NewEncoder(rw).Encode(v); err != nil {
		log.Println("Could not encode result for ", id)
		rw.WriteHeader(500)
	}
}

func generateData(d time.Time) {

	log.Println("generating data...")

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
	go stageRead(*dataPath, stage1, &readGroup, counter1)

	// Stage two progates an object forward one day and then passes it on.
	for i := 0; i < processors; i++ {
		fanGroup.Add(1)
		go stageMeanMotion(stage1, stage2, d, &fanGroup, counter2)
	}

	for i := 0; i < numFiles; i++ {
		// The final stage converts the orbit information into a position in space and then writes it to a file.
		complete.Add(1)
		go stageOutput(filename(i), stage2, &complete, counter3)
	}

	readGroup.Wait()

	fanGroup.Wait()
	close(stage2)

	complete.Wait()
	log.Println("done generating data")

}

// stageRead opens a file using the gompcreader project and reads out orbital information.
func stageRead(inputFile string, output chan *orbcore.Orbit, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	mpcReader, err := gompcreader.NewMpcReader(inputFile)
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

		output <- orb

		counter.Incr(1)
		count++
		result, err = mpcReader.ReadEntry()
	}
	if err != io.EOF {
		log.Fatal("error reading", err)
	}
}

func stageMeanMotion(in chan *orbcore.Orbit, output chan *orbcore.Position, targetDate time.Time, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()
	for orb := range in {
		orb = orbcore.MeanMotionToDate(orb, targetDate)
		output <- orbcore.OrbitToPosition(orb)

		counter.Incr(1)
	}
}

func stageOutput(outputPath string, in chan *orbcore.Position, wg *sync.WaitGroup, counter *ratecounter.RateCounter) {
	defer wg.Done()

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("error creating outputfile ", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Could not close file ", outputPath, err)
		}
	}()

	w := bufio.NewWriterSize(f, 64*1024)
	defer func() {
		if err := w.Flush(); err != nil {
			log.Fatal("Could not flush file ", outputPath, err)
		}
	}()

	for orb := range in {
		if _, err := fmt.Fprintf(f, "%s,%v,%v,%v\n", orb.ID, orb.X, orb.Y, orb.Z); err != nil {
			log.Fatal("Could not write", err)
		}
		counter.Incr(1)
	}

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

func prepObjectData(inputPath string) error {
	mpcReader, err := gompcreader.NewMpcReader(inputPath)
	if err != nil {
		log.Fatal("error creating mpcReader ", err)
	}
	defer mpcReader.Close()

	asteroidData = make(map[string]objectData)

	result, err := mpcReader.ReadEntry()
	for err == nil {
		ob := objectData{
			ID:                          result.ID,
			Epoch:                       result.Epoch,
			MeanAnomalyEpoch:            orbconvert.DegToRad(result.MeanAnomalyEpoch),
			ArgumentOfPerihelion:        orbconvert.DegToRad(result.ArgumentOfPerihelion),
			LongitudeOfTheAscendingNode: orbconvert.DegToRad(result.LongitudeOfTheAscendingNode),
			InclinationToTheEcliptic:    orbconvert.DegToRad(result.InclinationToTheEcliptic),
			OrbitalEccentricity:         result.OrbitalEccentricity,
			MeanDailyMotion:             result.MeanDailyMotion,
			SemimajorAxis:               orbconvert.AuToKm(result.SemimajorAxis),
			Orbit:                       nil,
		}

		asteroidData[ob.ID] = ob
		result, err = mpcReader.ReadEntry()
	}
	if err != io.EOF {
		return err
	}
	return nil
}

type objectData struct {
	ID                          string
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

func (o *objectData) toOrbit() *orbcore.Orbit {
	return &orbcore.Orbit{
		ID:                          o.ID,
		ParentGrav:                  orbdata.SunGrav,
		Epoch:                       o.Epoch,
		MeanAnomalyEpoch:            o.MeanAnomalyEpoch,
		ArgumentOfPerihelion:        o.ArgumentOfPerihelion,
		LongitudeOfTheAscendingNode: o.LongitudeOfTheAscendingNode,
		InclinationToTheEcliptic:    o.InclinationToTheEcliptic,
		OrbitalEccentricity:         o.OrbitalEccentricity,
		MeanDailyMotion:             o.MeanDailyMotion,
		SemimajorAxis:               o.SemimajorAxis,
	}
}

type point struct {
	X float64
	Y float64
	Z float64
}
