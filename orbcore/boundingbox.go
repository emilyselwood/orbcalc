package orbcore

import (
	"fmt"
	"math"
	"time"
)

/*
BoundingBox defines a four dimensional bounding box in space and time.
*/
type BoundingBox struct {
	MinX    float64
	MinY    float64
	MinZ    float64
	MinTime time.Time
	MaxX    float64
	MaxY    float64
	MaxZ    float64
	MaxTime time.Time
}

func (bb *BoundingBox) String() string {
	return fmt.Sprintf(
		"(%v,%v,%v,%v)x(%v,%v,%v,%v)",
		bb.MinX, bb.MinY, bb.MinZ, bb.MinTime,
		bb.MaxX, bb.MaxY, bb.MaxZ, bb.MaxTime,
	)
}

/*
Contains returns true if the provided position is inside the box.
*/
func (bb *BoundingBox) Contains(pos *Position) bool {
	return bb.MinX <= pos.X && pos.X <= bb.MaxX &&
		bb.MinY <= pos.Y && pos.Y <= bb.MaxY &&
		bb.MinZ <= pos.Z && pos.Z <= bb.MaxZ &&
		(bb.MinTime.Before(pos.Epoch) || bb.MinTime.Equal(pos.Epoch)) &&
		(pos.Epoch.Before(bb.MaxTime) || pos.Epoch.Equal(bb.MaxTime))
}

/*
Overlaps returns true if the other bounding box overlaps this one.
*/
func (bb *BoundingBox) Overlaps(other *BoundingBox) bool {
	// For each of the corners of this bounding box. are they inside the other bounding box?
	for _, c := range bb.Corners() {
		if other.Contains(&c) {
			return true
		}
	}

	// If not are any of the corners of the other box inside this box?
	for _, c := range other.Corners() {
		if bb.Contains(&c) {
			return true
		}
	}

	return false
}

/*
Center returns the point in the center of this bounding box.
*/
func (bb *BoundingBox) Center() Position {
	return Position{
		ID:    "center",
		Epoch: splitTime(bb.MinTime, bb.MaxTime),
		X:     splitFloat64(bb.MinX, bb.MaxX),
		Y:     splitFloat64(bb.MinY, bb.MaxY),
		Z:     splitFloat64(bb.MinZ, bb.MaxZ),
	}
}

/*
Corners returns the 16 corners of our 4 dimensional bounding box.
*/
func (bb *BoundingBox) Corners() [16]Position {
	var result [16]Position

	for i := 0; i < 16; i++ {
		result[i].X = pickSide(i&0x1, bb.MinX, bb.MaxX)
		result[i].Y = pickSide(i&0x2, bb.MinY, bb.MaxY)
		result[i].Z = pickSide(i&0x4, bb.MinZ, bb.MaxZ)
		if i&0x8 == 0 {
			result[i].Epoch = bb.MinTime
		} else {
			result[i].Epoch = bb.MaxTime
		}
	}
	return result
}

func pickSide(side int, min, max float64) float64 {
	if side == 0 {
		return min
	} else {
		return max
	}
}

type BoundingBoxSplitter struct {
	Box    *BoundingBox
	Center Position
	i      int
}

func NewBoundingBoxSplitter(box *BoundingBox) BoundingBoxSplitter {
	return BoundingBoxSplitter{
		Box:    box,
		Center: box.Center(),
		i:      0,
	}
}

func (bs BoundingBoxSplitter) Next() *BoundingBox {
	result := bs.At(bs.i)
	bs.i = bs.i + 1
	return result
}

func (bs BoundingBoxSplitter) HasNext() bool {
	return bs.i < 16
}

func (bs BoundingBoxSplitter) At(i int) *BoundingBox {
	var result BoundingBox
	if i >= 16 {
		return &result
	}

	result.MinX, result.MaxX = pickBoxEdge(i&0x1, bs.Box.MinX, bs.Center.X, bs.Box.MaxX)
	result.MinY, result.MaxY = pickBoxEdge(i&0x2, bs.Box.MinY, bs.Center.Y, bs.Box.MaxY)
	result.MinZ, result.MaxZ = pickBoxEdge(i&0x4, bs.Box.MinZ, bs.Center.Z, bs.Box.MaxZ)

	if i&0x8 == 0 {
		result.MinTime = bs.Box.MinTime
		result.MaxTime = bs.Center.Epoch
	} else {
		result.MinTime = bs.Center.Epoch
		result.MaxTime = bs.Box.MaxTime
	}

	return &result
}

func pickBoxEdge(side int, min, mid, max float64) (float64, float64) {
	if side == 0 {
		return min, mid
	} else {
		return mid, max
	}
}

/*
PositionsToBoundingBox creates a bounding box around a set of positions.
*/
func PositionsToBoundingBox(positions []*Position) *BoundingBox {

	minX, minY, minZ := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	maxX, maxY, maxZ := -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64

	minTime := time.Unix(1<<63-62135596801, 999999999)
	maxTime := time.Unix(-62135596801, -999999999)

	for _, p := range positions {

		minX = math.Min(minX, p.X)
		minY = math.Min(minY, p.Y)
		minZ = math.Min(minZ, p.Z)

		maxX = math.Max(maxX, p.X)
		maxY = math.Max(maxY, p.Y)
		maxZ = math.Max(maxZ, p.Z)

		if p.Epoch.Before(minTime) {
			minTime = p.Epoch
		}

		if p.Epoch.After(maxTime) {
			maxTime = p.Epoch
		}
	}

	return &BoundingBox{
		MinX: minX, MaxX: maxX,
		MinY: minY, MaxY: maxY,
		MinZ: minZ, MaxZ: maxZ,
		MinTime: minTime, MaxTime: maxTime,
	}
}

/*
splitFloat64 finds the mid point between two floats
*/
func splitFloat64(min, max float64) float64 {
	return ((max - min) / 2) + min
}

func splitTime(min, max time.Time) time.Time {
	return min.Add(time.Duration(int64(max.Sub(min)) / 2))
}
