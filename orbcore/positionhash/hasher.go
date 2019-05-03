package positionhash

import (
	"fmt"
	"github.com/wselwood/orbcalc/orbcore"
	"strings"
	"time"
)

/*
Hasher defines a way to create spacial temporal hashes.
*/
type Hasher interface {
	Hash(pos *orbcore.Position) (string, error)
	Box(hash string) (*orbcore.BoundingBox, error)
}

/*
HexHash is a Hasher that uses 16 buckets per level. This splits each dimension in half every go. Ending up with a binary
tree across four dimensions.

The idea here is like a geohash but across more dimensions so we can define a box of space and time and easily match
positions that are in the box or not.

An instance of HexHasher is not thread safe.
 */
type HexHasher struct {
	Space *orbcore.BoundingBox
	Depth int
	boxBuffer [16]orbcore.BoundingBox
}

func (hh *HexHasher) Hash(pos *orbcore.Position) (string, error) {
	builder, err := hh.generateHexHash(pos, *(hh.Space), &strings.Builder{})
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (hh *HexHasher) Box(hash string) (orbcore.BoundingBox, error) {
	return hh.findBox(hash, *(hh.Space))
}

func (hh *HexHasher) generateHexHash(pos *orbcore.Position, box orbcore.BoundingBox, result *strings.Builder) (*strings.Builder, error) {
	if !box.Contains(pos) {
		return result, fmt.Errorf("position is not valid for this hasher")
	}

	if result.Len() == hh.Depth {
		return result, nil
	}

	const values string = "0123456789ABCDEF"
	splits := splitBox(box, hh.boxBuffer)
	for i, b := range splits {
		if b.Contains(pos) {
			result.WriteRune(rune(values[i]))
			return hh.generateHexHash(pos, b, result)
		}
	}
	return result, fmt.Errorf("could not find sub bounding box to select from %v for point %v", box, pos)
}

func (hh *HexHasher) findBox(hash string, parent orbcore.BoundingBox) (orbcore.BoundingBox, error) {
	if hash == "" {
		return parent, nil
	}

	index := int(hexToDec(hash[0]))
	if index < 0 || index >= 16 {
		return parent, fmt.Errorf("unknown character in hash, 0-9A-F are valid")
	}
	splits := splitBox(parent, hh.boxBuffer)

	return hh.findBox(hash[1:], splits[index])
}

/*
splitBox cuts a bounding box in two along all of its dimensions
*/
func splitBox(box orbcore.BoundingBox, result [16]orbcore.BoundingBox) [16]orbcore.BoundingBox {

	minX, midX, maxX := splitFloat64(box.MinX, box.MaxX)
	minY, midY, maxY := splitFloat64(box.MinY, box.MaxY)
	minZ, midZ, maxZ := splitFloat64(box.MinZ, box.MaxZ)
	minTime, midTime, maxTime := splitTime(box.MinTime, box.MaxTime)

	// We need to have each combination of min and max boxes for the four dimensions.
	// use the binary encoding of an int to achieve this.
	for i := 0; i < 16; i++ {
		result[i].MinX, result[i].MaxX = pickSide(i&0x1, minX, midX, maxX)
		result[i].MinY, result[i].MaxY = pickSide(i&0x2, minY, midY, maxY)
		result[i].MinZ, result[i].MaxZ = pickSide(i&0x4, minZ, midZ, maxZ)

		if i&0x8 == 0 {
			result[i].MinTime = minTime
			result[i].MaxTime = midTime
		} else {
			result[i].MinTime = midTime
			result[i].MaxTime = maxTime
		}
	}

	return result
}

func pickSide(side int, min, mid, max float64) (float64, float64) {
	if side == 0 {
		return min, mid
	} else {
		return mid, max
	}
}

/*
splitFloat64 finds the mid point between two floats
*/
func splitFloat64(min, max float64) (float64, float64, float64) {
	return min, ((max - min) / 2) + min, max
}

func splitTime(min, max time.Time) (time.Time, time.Time, time.Time) {
	return min, min.Add(time.Duration(int64(max.Sub(min)) / 2)), max
}

func hexToDec(c uint8) uint8 {
	if c >= '0' && c <= '9' {
		return c - '0'
	} else {
		return (c - 'A') + 10
	}
}