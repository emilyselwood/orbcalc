/*
Package orbcore provides functions to work out orbital details for an asronomical object.
*/
package orbcore

import (
	"fmt"
	"math"
	"time"

	"github.com/wselwood/orbcalc/orbdata"
	"gonum.org/v1/gonum/mat"
)

/*
Orbit holds required information for orbit calculations
*/
type Orbit struct {
	ID                          string
	AbsoluteMagnitude           float64
	Slope                       float64
	Epoch                       time.Time
	MeanAnomalyEpoch            float64 // nu
	ArgumentOfPerihelion        float64 // w argp
	LongitudeOfTheAscendingNode float64 // omega raan
	InclinationToTheEcliptic    float64 // i inc
	OrbitalEccentricity         float64 // e ecc
	MeanDailyMotion             float64
	SemimajorAxis               float64 // a p
}

/*
Clone makes a copy of this orbit object
*/
func (o *Orbit) Clone() *Orbit {
	return &Orbit{
		ID:                          o.ID,
		AbsoluteMagnitude:           o.AbsoluteMagnitude,
		Slope:                       o.Slope,
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

func (o *Orbit) String() string {
	return fmt.Sprintf("id: \"%v\" nu: %v w: %v omega: %v i: %v e: %v a: %v",
		o.ID,
		o.MeanAnomalyEpoch,
		o.ArgumentOfPerihelion,
		o.LongitudeOfTheAscendingNode,
		o.InclinationToTheEcliptic,
		o.OrbitalEccentricity,
		o.SemimajorAxis,
	)
}

/*
VectorToHelocentric converts to helocentric reference frame from the default used by the MPC
*/
func VectorToHelocentric(r mat.Vector, v mat.Vector) (mat.Vector, mat.Vector) {

	// TODO: Implement this

	return r, v
}

/*
OrbitToVector creates a vector representation from a MinorPlanet object
*/
func OrbitToVector(orbit *Orbit) (mat.Vector, mat.Vector) {

	r, v := OrbitToVecPerifocal(orbit)

	/*
		TODO: Fix this and put it back in some day.
		rot := QuickerRotationMatrixForOrbit(orbit.ArgumentOfPerihelion, orbit.InclinationToTheEcliptic, orbit.LongitudeOfTheAscendingNode)

		r.MulVec(rot, r)
		v.MulVec(rot, v)
	*/
	r = Rotate(r, orbit.ArgumentOfPerihelion, AxisZ)
	r = Rotate(r, orbit.InclinationToTheEcliptic, AxisX)
	r = Rotate(r, orbit.LongitudeOfTheAscendingNode, AxisZ)

	v = Rotate(v, orbit.ArgumentOfPerihelion, AxisZ)
	v = Rotate(v, orbit.InclinationToTheEcliptic, AxisX)
	v = Rotate(v, orbit.LongitudeOfTheAscendingNode, AxisZ)

	return r, v
}

// TODO: VectorToOrbit

/*
OrbitToVecPerifocal converts a MinorPlanet object into r and v vectors in the perifocal frame
*/
func OrbitToVecPerifocal(orbit *Orbit) (*mat.VecDense, *mat.VecDense) {

	a := orbit.SemimajorAxis * (1 - math.Pow(orbit.OrbitalEccentricity, 2))
	cosNu := math.Cos(orbit.MeanAnomalyEpoch)
	sinNu := math.Sin(orbit.MeanAnomalyEpoch)

	r := mat.NewVecDense(3, []float64{cosNu, sinNu, 0 * orbit.MeanAnomalyEpoch})

	rMult := a / (1 + orbit.OrbitalEccentricity*cosNu)
	r.ScaleVec(rMult, r)

	v := mat.NewVecDense(3, []float64{-sinNu, orbit.OrbitalEccentricity + cosNu, 0})

	vMult := math.Sqrt(orbdata.SunGrav / a)
	v.ScaleVec(vMult, v)

	return r, v
}

/*
OrbitalPeriod returns the time taken for a complete orbit.
*/
func OrbitalPeriod(orbit *Orbit, parent float64) time.Duration {
	t := 2 * math.Pi * math.Sqrt(math.Pow(orbit.SemimajorAxis, 3)/parent)
	return time.Duration(t) * time.Second
}
