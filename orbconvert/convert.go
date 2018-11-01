package orbconvert

import (
	"math"

	"github.com/wselwood/gompcreader"
	"github.com/wselwood/orbcalc/orbcore"
	"github.com/wselwood/orbcalc/orbdata"
)

/*
ConvertFromMinorPlanet converts the columns in a minor planet into the requred types for the rest of this to work
*/
func ConvertFromMinorPlanet(mpc *gompcreader.MinorPlanet) *orbcore.Orbit {
	var result orbcore.Orbit

	// TODO: convert more types as needed
	result.ID = mpc.ID
	result.AbsoluteMagnitude = mpc.AbsoluteMagnitude
	result.Slope = mpc.Slope
	result.Epoch = mpc.Epoch
	result.MeanAnomalyEpoch = DegToRad(mpc.MeanAnomalyEpoch)
	result.ArgumentOfPerihelion = DegToRad(mpc.ArgumentOfPerihelion)
	result.LongitudeOfTheAscendingNode = DegToRad(mpc.LongitudeOfTheAscendingNode)
	result.InclinationToTheEcliptic = DegToRad(mpc.InclinationToTheEcliptic)
	result.OrbitalEccentricity = mpc.OrbitalEccentricity
	result.MeanDailyMotion = mpc.MeanDailyMotion
	result.SemimajorAxis = AuToKm(mpc.SemimajorAxis)

	return &result
}

const toRad = math.Pi / 180.0

/*
DegToRad converts degrees to radians
*/
func DegToRad(in float64) float64 {
	return in * toRad
}

/*
AuToKm converts a distance in AU into Km
*/
func AuToKm(in float64) float64 {
	return in * orbdata.AU
}

/*
KmToAu converts a distance in Km to AU
*/
func KmToAu(in float64) float64 {
	return in / orbdata.AU
}
