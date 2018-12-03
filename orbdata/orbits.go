package orbdata

import (
	"github.com/wselwood/orbcalc/orbcore"
)

// This file will contain orbital information for standard objects. Major planets, moons and so on.

// EarthOrbit defines the standard earth orbit.
var EarthOrbit = orbcore.Orbit{
	ID: "Earth",
	MeanAnomalyEpoch: 6.039693392708146, // rad
	ArgumentOfPerihelion: 1.4877567222443007, // rad
	LongitudeOfTheAscendingNode: 8.219803446009808e-05, // rad
	InclinationToTheEcliptic: 0.40898784995301246, // rad
	OrbitalEccentricity: 0.023506053256160484,
	SemimajorAxis: 151869811.1699976, // km
}
