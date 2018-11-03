package orbcore

import (
	"testing"

	"github.com/wselwood/orbcalc/orbdata"
)

func TestLoopingProblem(t *testing.T) {

	orb := Orbit{
		ID:                          "1996 PW",
		MeanAnomalyEpoch:            0.03539440456581901,
		ArgumentOfPerihelion:        3.169512336568096,
		LongitudeOfTheAscendingNode: 2.519967809619083,
		InclinationToTheEcliptic:    0.5228416517687837,
		OrbitalEccentricity:         0.9901593,
		SemimajorAxis:               3.79035922723884e+10,
	}

	_ = MeanMotion(orbdata.SunGrav, &orb, 1*24*60*60)

}
