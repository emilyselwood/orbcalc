/*
This tool takes several defined orbits, propagates them forward in time and writes the results to a text file.

This is used so that we can compare the results with other astrodynamics packages. Tests that should generate the
same output are found in sub folders from here.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
)

var testObjects = []*orbcore.Orbit{
	{
		ID:                          "1996 PW",
		MeanAnomalyEpoch:            0.03539440456581901,
		ArgumentOfPerihelion:        3.169512336568096,
		LongitudeOfTheAscendingNode: 2.519967809619083,
		InclinationToTheEcliptic:    0.5228416517687837,
		OrbitalEccentricity:         0.9901593,
		SemimajorAxis:               3.79035922723884e+10,
	},
}

func main() {

	output := flag.String("out", "./", "output directory")

	flag.Parse()

	if *output == "" {
		log.Fatal("output directory needs to be defined.")
	}

	for i, o := range testObjects {
		if err := processOrbit(o, *output); err != nil {
			log.Fatalf("could not process entry %d got error %v", i, err)
		}
	}

}

func processOrbit(orb *orbcore.Orbit, outDir string) error {

	f, err := os.Create(filepath.Join(outDir, orb.ID+".csv"))
	if err != nil {
		return err
	}

	defer f.Close()

	for i := 1; i <= 365; i++ { // loop for a year of days
		seconds := int64(i * (24 * 60 * 60)) // seconds per day
		updated := orbcore.MeanMotion(orbdata.SunGrav, orb, seconds)
		r, _ := orbcore.OrbitToVector(updated)
		fmt.Fprintf(f, "%v,%v,%v,%v,%v\n", orb.ID, i, r.AtVec(0), r.AtVec(1), r.AtVec(2))
	}

	return nil
}
