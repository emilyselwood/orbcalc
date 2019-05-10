package orbdata

import (
	"github.com/wselwood/orbcalc/orbcore"
)

// This file will contain orbital information for standard objects. Major planets, moons and so on.

// MercuryOrbit defines the standard mercury orbit
var MercuryOrbit = orbcore.Orbit{
	ID:                          "Mercury",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            0.7363828677023899,  // rad
	ArgumentOfPerihelion:        1.290398137330985,   // rad
	LongitudeOfTheAscendingNode: 0.19016162418731905, // rad
	InclinationToTheEcliptic:    0.122258,            // rad
	OrbitalEccentricity:         0.2161872518335417,
	SemimajorAxis:               5.7909176e7, // km
}

// VenusOrbit defines the standard mercury orbit
var VenusOrbit = orbcore.Orbit{
	ID:                          "Venus",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            6.024347789858294,   // rad
	ArgumentOfPerihelion:        1.8790979389622697,  // rad
	LongitudeOfTheAscendingNode: 0.13963804205942293, // rad
	InclinationToTheEcliptic:    0.0592489,           // rad
	OrbitalEccentricity:         0.017361719534212148,
	SemimajorAxis:               1.0820893e8, // km
}

// EarthOrbit defines the standard earth orbit.
var EarthOrbit = orbcore.Orbit{
	ID:                          "Earth",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            6.039693392708146,     // rad
	ArgumentOfPerihelion:        1.4877567222443007,    // rad
	LongitudeOfTheAscendingNode: 8.219803446009808e-05, // rad
	InclinationToTheEcliptic:    0.9e-07,               // rad
	OrbitalEccentricity:         0.023506053256160484,
	SemimajorAxis:               1.49597887e8, // km
}

// MarsOrbit defines the standard mars orbit.
var MarsOrbit = orbcore.Orbit{
	ID:                          "Mars",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            0.9016227920497925,   // rad
	ArgumentOfPerihelion:        5.804221558977953,    // rad
	LongitudeOfTheAscendingNode: 0.059136325715984754, // rad
	InclinationToTheEcliptic:    0.0322992,            // rad
	OrbitalEccentricity:         0.09853112210172534,
	SemimajorAxis:               2.27936637e8, // km
}

// JupiterOrbit defines the standard Jupiter orbit.
var JupiterOrbit = orbcore.Orbit{
	ID:                          "Jupiter",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            3.986624571747394,    // rad
	ArgumentOfPerihelion:        0.22894709895829354,  // rad
	LongitudeOfTheAscendingNode: 0.056682739190454204, // rad
	InclinationToTheEcliptic:    0.0227818,            // rad
	OrbitalEccentricity:         0.05041232826440195,
	SemimajorAxis:               7.78412027e8, // km
}

// SaturnOrbit defines the standard Saturn orbit.
var SaturnOrbit = orbcore.Orbit{
	ID:                          "Saturn",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            3.2720797523951766,  // rad
	ArgumentOfPerihelion:        1.5276434137035415,  // rad
	LongitudeOfTheAscendingNode: 0.10399170848152173, // rad
	InclinationToTheEcliptic:    0.043362,            // rad
	OrbitalEccentricity:         0.05853326249640754,
	SemimajorAxis:               1.42672541e9, // km
}

// UranusOrbit defines the standard Uranus orbit.
var UranusOrbit = orbcore.Orbit{
	ID:                          "Uranus",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            3.8644829632802806,  // rad
	ArgumentOfPerihelion:        2.926548412800625,   // rad
	LongitudeOfTheAscendingNode: 0.03235322856941487, // rad
	InclinationToTheEcliptic:    0.013437,            // rad
	OrbitalEccentricity:         0.044645557888114,
	SemimajorAxis:               2.87097222e9, // km
}

// NeptuneOrbit defines the standard Neptune orbit.
var NeptuneOrbit = orbcore.Orbit{
	ID:                          "Neptune",
	ParentGrav:                  SunGrav,
	Epoch:                       J2000,
	MeanAnomalyEpoch:            5.100969108525634,    // rad
	ArgumentOfPerihelion:        0.8712884041923264,   // rad
	LongitudeOfTheAscendingNode: 0.060720496894987035, // rad
	InclinationToTheEcliptic:    0.0308778,            // rad
	OrbitalEccentricity:         0.011600603763700122,
	SemimajorAxis:               4.49825291e9, // km
}

// SolarSystem is a collection of major bodies in the solar system
var SolarSystem = []orbcore.Orbit{
	MercuryOrbit,
	VenusOrbit,
	EarthOrbit,
	MarsOrbit,
	JupiterOrbit,
	SaturnOrbit,
	UranusOrbit,
	NeptuneOrbit,
}

// InnerSolarSystem is a collection of major bodies in the inner solar system
var InnerSolarSystem = []orbcore.Orbit{
	MercuryOrbit,
	VenusOrbit,
	EarthOrbit,
	MarsOrbit,
}

// OuterSolarSystem is a collection of major bodies in the outer solar system
var OuterSolarSystem = []orbcore.Orbit{
	JupiterOrbit,
	SaturnOrbit,
	UranusOrbit,
	NeptuneOrbit,
}
