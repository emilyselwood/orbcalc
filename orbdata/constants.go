/*
Package orbdata contains useful contstants and defined orbits of the solar system.

In here you will find Gravitational constants of the major planets
*/
package orbdata

import (
	"time"
)

/*
SunGrav IAU2009 Heliocentric gravitational constant (km^3)/(s^-2)
*/
const SunGrav = 132712442099.00002

// MercuryGrav gravitational constant for Mercury centric orbits (km^3)/(s^-2)
const MercuryGrav = 22032.09

// VenusGrav gravitational constant for Venus centric orbits (km^3)/(s^-2)
const VenusGrav = 324858.592

// EarthGrav gravitational constant for Earth centric orbits (km^3)/(s^-2)
const EarthGrav = 398600.44180000003

// MoonGrav gravitational constant for Moon centric orbits (km^3)/(s^-2)
const MoonGrav = 4902.79981

// MarsGrav gravitational constant for Mars centric orbits (km^3)/(s^-2)
const MarsGrav = 42828.3744

// JupiterGrav gravitational constant for Jupiter centric orbits (km^3)/(s^-2)
const JupiterGrav = 126712762.53

// SaturnGrav gravitational constant for Saturn centric orbits (km^3)/(s^-2)
const SaturnGrav = 37931207.7

// UranusGrav gravitational constant for Uranus centric orbits (km^3)/(s^-2)
const UranusGrav = 5793939.300000001

// NeptuneGrav gravitational constant for Neptune centric orbits (km^3)/(s^-2)
const NeptuneGrav = 6836527.100580397

/*
AU represents the length of an Astronomical Unit in KiloMeters
*/
const AU = 149598000

/*
J2000 is the base time epoch of a lot of astronomical times.
*/
var J2000 = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
