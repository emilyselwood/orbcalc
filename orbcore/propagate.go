package orbcore

import (
	"fmt"
	"log"
	"math"
	"time"
)

/*
MeanMotionStepped calculates a mean motion value for [count] [timeStep]s and returns a list of orbits.
Note: The first entry in the returned list will always be the starting orbit object.
*/
func MeanMotionStepped(parent float64, orbit *Orbit, timeStep time.Duration, count int64) []*Orbit {
	result := make([]*Orbit, count+1)
	result[0] = orbit
	var i int64
	for i = 1; i <= count; i++ {
		offset := timeStep * time.Duration(i)
		result[i] = MeanMotion(parent, orbit, offset)
	}
	return result
}

/*
MeanMotionSteppedChannel works like MeanMotionStepped except it puts the results down a channel rather than returning a list
*/
func MeanMotionSteppedChannel(parent float64, orbit *Orbit, timeStep time.Duration, count int64, output chan *Orbit) {
	var i int64
	for i = 1; i <= count; i++ {
		offset := timeStep * time.Duration(i)
		output <- MeanMotion(parent, orbit, offset)
	}
}

/*
MeanMotionFullOrbit will calculate a number of entries for a full orbit, divided into [count] steps
*/
func MeanMotionFullOrbit(parent float64, orbit *Orbit, count int64) []*Orbit {
	orbitalPeriod := OrbitalPeriod(orbit, parent)
	step := orbitalPeriod / time.Duration(count)
	return MeanMotionStepped(parent, orbit, step, count)
}

/*
MeanMotion uses the mean motion method to propgate [orbit] through [t] seconds around [parent].
*/
func MeanMotion(parent float64, orbit *Orbit, t time.Duration) *Orbit {
	p := orbit.SemimajorAxis * (1 - math.Pow(orbit.OrbitalEccentricity, 2))
	m0 := createM0(orbit)
	var newMeanAnomalyEpoch float64
	if math.Abs(orbit.OrbitalEccentricity-1) > delta {
		a := p / (1 - math.Pow(orbit.OrbitalEccentricity, 2))
		m := m0 + float64(t.Seconds())*math.Sqrt(parent/math.Abs(math.Pow(a, 3)))
		newMeanAnomalyEpoch = mtoMeanAnomaly(m, orbit)
	} else {
		q := p * math.Abs(1.0-orbit.OrbitalEccentricity) / math.Abs(1.0-math.Pow(orbit.OrbitalEccentricity, 2))
		m := m0 + float64(t.Seconds())*math.Sqrt(parent/2.0/math.Pow(q, 3))
		newMeanAnomalyEpoch = mtoMeanAnomaly(m, orbit)
	}

	r := orbit.Clone()
	r.MeanAnomalyEpoch = newMeanAnomalyEpoch
	r.Epoch = r.Epoch.Add(t)

	return r
}

const delta = 1e-3
const tolerance = 1e-16

func createM0(orbit *Orbit) float64 {
	if orbit.OrbitalEccentricity > 1+delta {
		f := math.Log((math.Sqrt(orbit.OrbitalEccentricity+1) + math.Sqrt(orbit.OrbitalEccentricity-1)*math.Tan(orbit.MeanAnomalyEpoch/2)) /
			(math.Sqrt(orbit.OrbitalEccentricity+1) - math.Sqrt(orbit.OrbitalEccentricity-1)*math.Tan(orbit.MeanAnomalyEpoch/2)))
		return -f + orbit.OrbitalEccentricity*math.Sinh(f)
	} else if orbit.OrbitalEccentricity < 1-delta {
		e := 2 * math.Atan(math.Sqrt((1-orbit.OrbitalEccentricity)/(1+orbit.OrbitalEccentricity))*math.Tan(orbit.MeanAnomalyEpoch/2))
		return e - orbit.OrbitalEccentricity*math.Sin(e)
	} else {
		d := math.Tan(orbit.MeanAnomalyEpoch / 2.0)
		return keplerParabolic(orbit.OrbitalEccentricity, d)
	}
}

func mtoMeanAnomaly(m float64, orbit *Orbit) float64 {
	if orbit.OrbitalEccentricity > 1+delta {
		f := newtonKeplerHyper(math.Asinh(m/orbit.OrbitalEccentricity), m, orbit.OrbitalEccentricity)
		return 2 * math.Atan((math.Exp(f)*math.Sqrt(orbit.OrbitalEccentricity+1)-math.Sqrt(orbit.OrbitalEccentricity+1))/
			(math.Exp(f)*math.Sqrt(orbit.OrbitalEccentricity-1)+math.Sqrt(orbit.OrbitalEccentricity-1)))
	} else if orbit.OrbitalEccentricity < 1-delta {
		e := newtonKepler(m, m, orbit.OrbitalEccentricity)
		return 2 * math.Atan(math.Sqrt((1+orbit.OrbitalEccentricity)/(1-orbit.OrbitalEccentricity))*math.Tan(e/2))
	} else {
		b := 3.0 * m / 2.0
		a := (b + (1 + math.Pow(math.Pow(math.Pow(b, 2), 0.5), (2.0/3.0))))

		guess := 2 * a * b / (1 + a + math.Pow(a, 2))
		d := newtonKeplerParabolic(guess, m, orbit.OrbitalEccentricity)
		return 2.0 * math.Atan(d)
	}
}

// calculate kepler parabolic value
func keplerParabolic(orbitalEccentricity float64, d float64) float64 {

	x := (orbitalEccentricity - 1) / (orbitalEccentricity + 1) * math.Pow(d, 2)

	done := false
	s := 0.0
	k := 0.0
	fmt.Println(orbitalEccentricity, ",", d)
	for !done {
		term := (orbitalEccentricity - 1.0/(2.0*k+3.0)) * math.Pow(x, k)
		done = math.Abs(term) < tolerance
		s += term
		k = k + 1.0
	}

	return math.Sqrt(2.0/(1.0+orbitalEccentricity))*d + math.Sqrt(2.0/math.Pow((1.0+orbitalEccentricity), 3))*math.Pow(d, 3)*s
}

func keplerParabolicPrime(orbitalEccentricity float64, d float64) float64 {

	x := (orbitalEccentricity - 1) / (orbitalEccentricity + 1) * math.Pow(d, 2)
	done := false
	s := 0.0
	k := 0

	for !done {
		term := (orbitalEccentricity - 1.0/(2.0*float64(k)+3.0)) * (2.0*float64(k) + 3.0) * math.Pow(x, float64(k))
		done = math.Abs(term) < tolerance
		s += term
		k++
	}

	return math.Sqrt(2.0/(1.0+orbitalEccentricity))*d + math.Sqrt(2.0/math.Pow((1.0+orbitalEccentricity), 3))*math.Pow(d, 2)*s
}

func keplerHyper(orbitalEccentricity float64, d float64) float64 {
	return -d + orbitalEccentricity*math.Sinh(d)
}

func keplerHyperPrime(orbitalEccentricity float64, f float64) float64 {
	return orbitalEccentricity*math.Cosh(f) - 1
}

func kepler(orbitalEccentricity float64, d float64) float64 {
	return d - orbitalEccentricity*math.Sin(d)
}

func keplerPrime(orbitalEccentricity float64, f float64) float64 {
	return 1 - orbitalEccentricity*math.Cos(f)
}

func newtonKeplerParabolic(x0 float64, m float64, orbitalEccentricity float64) float64 {
	return newton(x0, m, orbitalEccentricity, func(ecc, p, m float64) float64 {
		f := keplerParabolic(orbitalEccentricity, p) - m
		d := keplerParabolicPrime(orbitalEccentricity, p)
		return f / d
	})
}

func newtonKeplerHyper(x0 float64, m float64, orbitalEccentricity float64) float64 {
	return newton(x0, m, orbitalEccentricity, func(ecc, p, m float64) float64 {
		f := keplerHyper(orbitalEccentricity, p) - m
		d := keplerHyperPrime(orbitalEccentricity, p)
		return f / d
	})
}

func newtonKepler(x0 float64, m float64, orbitalEccentricity float64) float64 {
	return newton(x0, m, orbitalEccentricity, func(ecc, p, m float64) float64 {
		f := kepler(orbitalEccentricity, p) - m
		d := keplerPrime(orbitalEccentricity, p)
		return f / d
	})
}

func newton(x0 float64, m float64, orbitalEccentricity float64, factor func(float64, float64, float64) float64) float64 {
	p0 := 1.0 * x0
	for i := 0; i < 100; i++ { // max number of iterations to do

		p := p0 - factor(orbitalEccentricity, p0, m)
		if math.Abs(p-p0) < 1.48e-8 {
			return p
		}
		p0 = p
	}
	log.Println("newton did not converge")
	return 1
}
