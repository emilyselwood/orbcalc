package orbcore

import (
	"testing"
)

func TestLoopingProblem(t *testing.T) {
	// This object has an eccentricity that is very nearly 1 and thus parabolic.
	// This was causing an infinate loop in the mean motion calculation.
	// This test is here to make sure it does not come back.
	orb := Orbit{
		ID:                          "1996 PW",
		ParentGrav:                  132712442099.00002,
		MeanAnomalyEpoch:            0.03539440456581901,
		ArgumentOfPerihelion:        3.169512336568096,
		LongitudeOfTheAscendingNode: 2.519967809619083,
		InclinationToTheEcliptic:    0.5228416517687837,
		OrbitalEccentricity:         0.9901593,
		SemimajorAxis:               3.79035922723884e+10,
	}

	_ = MeanMotion(&orb, 1*24*60*60)
}
