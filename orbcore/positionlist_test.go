package orbcore

import (
	"testing"
	"time"
)

func TestBoundingBox_Contains(t *testing.T) {
	box := BoundingBox{
		MinX: -10, MaxX: 10,
		MinY: -10, MaxY: 10,
		MinZ: -10, MaxZ: 10,
		MinTime: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		MaxTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	box2 := BoundingBox{
		MinX: -10, MaxX: 0,
		MinY: -10, MaxY: 0,
		MinZ: -10, MaxZ: 0,
		MinTime: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		MaxTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	points := []struct {
		pos    Position
		box    BoundingBox
		inside bool
	}{
		{
			pos: Position{
				ID:    "MidPoint",
				Epoch: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
				X:     0,
				Y:     0,
				Z:     0,
			},
			box:    box,
			inside: true,
		},
		{
			pos: Position{
				ID:    "MaxEdge",
				Epoch: time.Date(2019, 5, 3, 13, 37, 12, 0, time.UTC),
				X:     10,
				Y:     10,
				Z:     10,
			},
			box:    box,
			inside: true,
		},
		{
			pos: Position{
				ID:    "MinEdge",
				Epoch: time.Date(2019, 5, 3, 13, 37, 12, 0, time.UTC),
				X:     -10,
				Y:     -10,
				Z:     -10,
			},
			box:    box,
			inside: true,
		},
		{
			pos: Position{
				ID:    "TopEdgeOffsetTime",
				Epoch: time.Date(2019, 5, 3, 13, 37, 12, 0, time.UTC),
				X:     0,
				Y:     0,
				Z:     0,
			},
			box:    box2,
			inside: true,
		},
		{
			pos: Position{
				ID:    "Outside",
				Epoch: time.Date(2019, 5, 3, 13, 37, 12, 0, time.UTC),
				X:     10,
				Y:     10,
				Z:     10,
			},
			box:    box2,
			inside: false,
		},
	}

	for _, p := range points {
		t.Run("BoundingBox_Contains "+p.pos.ID, func(t2 *testing.T) {
			if p.box.Contains(&p.pos) != p.inside {
				t2.Fatalf("Expected %v to be inside %v but it wasn't", p.pos, p.box)
			}
		})
	}

}
