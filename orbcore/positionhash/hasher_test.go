package positionhash

import (
	"github.com/wselwood/orbcalc/orbcore"
	"testing"
	"time"
)

func TestSplitBoundingBox(t *testing.T) {
	input := orbcore.BoundingBox{
		MinX: -10, MaxX: 10,
		MinY: -10, MaxY: 10,
		MinZ: -10, MaxZ: 10,
		MinTime: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		MaxTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	result := splitBox(input, [16]orbcore.BoundingBox{})
	if len(result) != 16 {
		t.Fatal("Wrong number of results returned.")
	}

	expectedOne := orbcore.BoundingBox{
		MinX: -10, MaxX: 0,
		MinY: -10, MaxY: 0,
		MinZ: -10, MaxZ: 0,
		MinTime: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		MaxTime: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	expectedEight := orbcore.BoundingBox{
		MinX: -10, MaxX: 0,
		MinY: -10, MaxY: 0,
		MinZ: -10, MaxZ: 0,
		MinTime: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		MaxTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if result[0] != expectedOne {
		t.Fatal("Did not get expected value for first entry in array.")
	}

	if result[8] != expectedEight {
		t.Fatalf("Did not get expected value for the forth entry in array. Got %v expected %v", result[4], expectedEight)
	}

}


func TestHexHasher_Hash(t *testing.T) {
	hasher := HexHasher{
		Space: &orbcore.BoundingBox{
			MinX: -10, MaxX: 10,
			MinY: -10, MaxY: 10,
			MinZ: -10, MaxZ: 10,
			MinTime: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			MaxTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Depth: 6,
	}

	input := orbcore.Position{
		ID: "wibble",
		Epoch: time.Date(2019, 5, 3, 13, 37, 12, 0, time.UTC),
		X: 0,
		Y: 0,
		Z: 0,
	}
	result, err := hasher.Hash(&input)
	if err != nil {
		t.Fatal(err)
	}
	if result != "8FF7FF" {
		t.Fatalf("expected 8FF7FF got '%v'", result)
	}
}


func TestHexHasher_Box(t *testing.T) {
	hasher := HexHasher{
		Space: &orbcore.BoundingBox{
			MinX: -10, MaxX: 10,
			MinY: -10, MaxY: 10,
			MinZ: -10, MaxZ: 10,
			MinTime: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			MaxTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Depth: 6,
	}

	box, err := hasher.Box("8FF7FF")
	if err != nil {
		t.Fatal(err)
	}
	expected := orbcore.BoundingBox{
		MinX : -0.3125,
		MinY : -0.3125,
		MinZ : -0.3125,
		MinTime: time.Date(2019, 3, 21, 16, 30, 0, 0, time.UTC),
		MaxX : 0,
		MaxY : 0,
		MaxZ : 0,
		MaxTime: time.Date(2019, 5, 17, 18, 0, 0, 0, time.UTC),
	}

	if expected != box {
		t.Fatalf("expected %v got %v", expected, box)
	}
}

func BenchmarkHash(b *testing.B) {
	hasher := HexHasher{
		Space: &orbcore.BoundingBox{
			MinX: -10, MaxX: 10,
			MinY: -10, MaxY: 10,
			MinZ: -10, MaxZ: 10,
			MinTime: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			MaxTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Depth: 6,
	}

	input := orbcore.Position{
		ID: "wibble",
		Epoch: time.Date(2019, 5, 3, 13, 37, 12, 0, time.UTC),
		X: 0,
		Y: 0,
		Z: 0,
	}

	for n:=0; n < b.N; n++ {
		result, err := hasher.Hash(&input)
		if err != nil {
			b.Fatal(err)
		}
		if result != "8FF7FF" {
			b.Fatalf("expected 8FF7FF got '%v'", result)
		}
	}
}