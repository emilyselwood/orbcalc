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
	"strings"
	"time"

	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
)

var testObjects = []*orbcore.Orbit{
	{
		ID:                          "1996 PW",                                   // Very very eliptical that can cause problems.
		Epoch:                       time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), // todo: real epoch times.
		MeanAnomalyEpoch:            0.03539440456581901,
		ArgumentOfPerihelion:        3.169512336568096,
		LongitudeOfTheAscendingNode: 2.519967809619083,
		InclinationToTheEcliptic:    0.5228416517687837,
		OrbitalEccentricity:         0.9901593,
		SemimajorAxis:               3.79035922723884e+10,
	},
	{
		ID:                          "1",                                         // Ceres
		Epoch:                       time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), // todo: real epoch times.
		MeanAnomalyEpoch:            6.147582300011738,
		ArgumentOfPerihelion:        1.2761023695175595,
		LongitudeOfTheAscendingNode: 1.4016725260132445,
		InclinationToTheEcliptic:    0.1848916288429445,
		OrbitalEccentricity:         0.0755347,
		SemimajorAxis:               4.1394459238740003e+08,
	},
	// TODO: More objects at least vesta and pluto as test cases
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

	f, err := os.Create(filepath.Join(outDir, cleanID(orb.ID)+".csv"))
	if err != nil {
		return err
	}

	defer f.Close()

	for _, r := range orbcore.MeanMotionStepped(orbdata.SunGrav, orb, time.Hour*24, 365) {
		p := orbcore.OrbitToPosition(r)
		fmt.Fprintln(f, p)
	}

	return nil
}

// Spaces are a pain in paths. Swap them for underscores.
func cleanID(id string) string {
	return strings.Replace(id, " ", "_", -1)
}
