package positionhash

import (
	"fmt"
	"github.com/emilyselwood/orbcalc/orbcore"
	"strings"
	"time"
)

/*
Hasher defines a way to create spacial temporal hashes.

Inspired by geohash, this extends to 4 dimensions.
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

An instance of HexHasher is *not* thread safe.
*/
type HexHasher struct {
	Space     *orbcore.BoundingBox
	Depth     int
	boxBuffer [16]orbcore.BoundingBox
	sb        strings.Builder
}

func (hh *HexHasher) Hash(pos *orbcore.Position) (string, error) {
	if !hh.Space.Contains(pos) {
		return "", fmt.Errorf("position is not valid for this hasher")
	}
	hh.sb.Reset()
	hh.sb.Grow(hh.Depth)
	err := hh.generateHexHash(pos, *(hh.Space))
	if err != nil {
		return "", err
	}
	return hh.sb.String(), nil
}

func (hh *HexHasher) Box(hash string) (orbcore.BoundingBox, error) {
	return hh.findBox(hash, *(hh.Space))
}


func (hh *HexHasher) generateHexHash(pos *orbcore.Position, box orbcore.BoundingBox) error {
	if hh.sb.Len() == hh.Depth {
		return nil
	}

	const hexValues string = "0123456789ABCDEF"
	splits := splitBox(box, hh.boxBuffer)
	for i := 0; i < 16; i ++ {
		b := splits[i]
		if b.Contains(pos) {
			hh.sb.WriteByte(hexValues[i])
			return hh.generateHexHash(pos, b)
		}
	}
	return fmt.Errorf("could not find sub bounding box to select from %v for point %v", box, pos)
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

	// unrolled array population to avoid branching at all.
	// yes this is ugly but it is damn quick
	// it replaces a loop that used a bit of an int to control which side of each field was used.
	result[0].MinX, result[0].MaxX = minX, midX
	result[0].MinY, result[0].MaxY = minY, midY
	result[0].MinZ, result[0].MaxZ = minZ, midZ
	result[0].MinTime, result[0].MaxTime = minTime, midTime

	result[1].MinX, result[1].MaxX = minX, maxX
	result[1].MinY, result[1].MaxY = midY, midY
	result[1].MinZ, result[1].MaxZ = minZ, midZ
	result[1].MinTime, result[1].MaxTime = minTime, midTime

	result[2].MinX, result[2].MaxX = minX, midX
	result[2].MinY, result[2].MaxY = midY, maxY
	result[2].MinZ, result[2].MaxZ = minZ, midZ
	result[2].MinTime, result[2].MaxTime = minTime, midTime

	result[3].MinX, result[3].MaxX = midX, maxX
	result[3].MinY, result[3].MaxY = midY, maxY
	result[3].MinZ, result[3].MaxZ = minZ, midZ
	result[3].MinTime, result[3].MaxTime = minTime, midTime

	result[4].MinX, result[4].MaxX = minX, midX
	result[4].MinY, result[4].MaxY = minY, midY
	result[4].MinZ, result[4].MaxZ = midZ, maxZ
	result[4].MinTime, result[4].MaxTime = minTime, midTime

	result[5].MinX, result[5].MaxX = midX, maxX
	result[5].MinY, result[5].MaxY = minY, midY
	result[5].MinZ, result[5].MaxZ = midZ, maxZ
	result[5].MinTime, result[5].MaxTime = minTime, midTime

	result[6].MinX, result[6].MaxX = minX, midX
	result[6].MinY, result[6].MaxY = midY, maxY
	result[6].MinZ, result[6].MaxZ = midZ, maxZ
	result[6].MinTime, result[6].MaxTime = minTime, midTime

	result[7].MinX, result[7].MaxX = midX, maxX
	result[7].MinY, result[7].MaxY = midY, maxY
	result[7].MinZ, result[7].MaxZ = midZ, maxZ
	result[7].MinTime, result[7].MaxTime = minTime, midTime

	result[8].MinX, result[8].MaxX = minX, midX
	result[8].MinY, result[8].MaxY = minY, midY
	result[8].MinZ, result[8].MaxZ = minZ, midZ
	result[8].MinTime, result[8].MaxTime = midTime, maxTime

	result[9].MinX, result[9].MaxX = minX, maxX
	result[9].MinY, result[9].MaxY = midY, midY
	result[9].MinZ, result[9].MaxZ = minZ, midZ
	result[9].MinTime, result[9].MaxTime = midTime, maxTime

	result[10].MinX, result[10].MaxX = minX, midX
	result[10].MinY, result[10].MaxY = midY, maxY
	result[10].MinZ, result[10].MaxZ = minZ, midZ
	result[10].MinTime, result[10].MaxTime = midTime, maxTime

	result[11].MinX, result[11].MaxX = midX, maxX
	result[11].MinY, result[11].MaxY = midY, maxY
	result[11].MinZ, result[11].MaxZ = minZ, midZ
	result[11].MinTime, result[11].MaxTime = midTime, maxTime

	result[12].MinX, result[12].MaxX = minX, midX
	result[12].MinY, result[12].MaxY = minY, midY
	result[12].MinZ, result[12].MaxZ = midZ, maxZ
	result[12].MinTime, result[12].MaxTime = midTime, maxTime

	result[13].MinX, result[13].MaxX = midX, maxX
	result[13].MinY, result[13].MaxY = minY, midY
	result[13].MinZ, result[13].MaxZ = midZ, maxZ
	result[13].MinTime, result[13].MaxTime = midTime, maxTime

	result[14].MinX, result[14].MaxX = minX, midX
	result[14].MinY, result[14].MaxY = midY, maxY
	result[14].MinZ, result[14].MaxZ = midZ, maxZ
	result[14].MinTime, result[14].MaxTime = midTime, maxTime

	result[15].MinX, result[15].MaxX = midX, maxX
	result[15].MinY, result[15].MaxY = midY, maxY
	result[15].MinZ, result[15].MaxZ = midZ, maxZ
	result[15].MinTime, result[15].MaxTime = midTime, maxTime

	return result
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
