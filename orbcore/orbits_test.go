package orbcore

import (
	"testing"
	"time"
)

func TestOrbitalPeriod(t *testing.T) {
	ceres := Orbit{
		ID:                          "1", // Ceres
		ParentGrav:                  132712442099.00002,
		Epoch:                       time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), // todo: real epoch times.
		MeanAnomalyEpoch:            6.147582300011738,
		ArgumentOfPerihelion:        1.2761023695175595,
		LongitudeOfTheAscendingNode: 1.4016725260132445,
		InclinationToTheEcliptic:    0.1848916288429445,
		OrbitalEccentricity:         0.0755347,
		SemimajorAxis:               4.1394459238740003e+08,
	}

	r := OrbitalPeriod(&ceres)

	expected, _ := time.ParseDuration("40349h11m39s")
	if r != expected {
		t.Errorf("Got %v expected %v", r, expected)
	}
}


func BenchmarkOrbitToVector(b *testing.B){

	ceres := Orbit{
		ID:                          "1", // Ceres
		ParentGrav:                  132712442099.00002,
		Epoch:                       time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), // todo: real epoch times.
		MeanAnomalyEpoch:            6.147582300011738,
		ArgumentOfPerihelion:        1.2761023695175595,
		LongitudeOfTheAscendingNode: 1.4016725260132445,
		InclinationToTheEcliptic:    0.1848916288429445,
		OrbitalEccentricity:         0.0755347,
		SemimajorAxis:               4.1394459238740003e+08,
	}

	for n := 0; n < b.N; n++ {
		r, v := OrbitToVector(&ceres)
		if r == nil || v == nil {
			b.Fatal("Got an invalid result")
		}
	}
}