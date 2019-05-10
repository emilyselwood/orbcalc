package orbcore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Position contains information about the location in space of an object
*/
type Position struct {
	ID    string
	Epoch time.Time
	X     float64
	Y     float64
	Z     float64
}

func (p *Position) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v", p.ID, p.Epoch.Format(time.RFC3339), p.X, p.Y, p.Z)
}

/*
OrbitToPosition converts an object object to its position on day 0
*/
func OrbitToPosition(orb *Orbit) *Position {
	r, _ := OrbitToVector(orb)
	return &Position{
		ID:    orb.ID,
		Epoch: orb.Epoch,
		X:     r.AtVec(0),
		Y:     r.AtVec(1),
		Z:     r.AtVec(2),
	}
}

/*
ReadPositionFile opens a path an loads in a CSV formatted list of Positions
*/
func ReadPositionFile(path string) (pos []*Position, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(c *os.File) {
		err = c.Close()
	}(file)

	return ReadPositions(file)
}

/*
ReadPositions takes a reader and parses a CSV formatted list of positions.
*/
func ReadPositions(input io.Reader) ([]*Position, error) {
	scanner := bufio.NewScanner(input)

	count := 0
	result := make([]*Position, 366)
	for scanner.Scan() {
		r, err := ParsePositionLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		if r != nil {
			if count >= len(result) {
				result = append(result, r)
				count++
			} else {
				result[count] = r
				count++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result[:count], nil
}

/*
ParsePositionLine takes a CSV formatted line representing a possition and returns the object or an error
*/
func ParsePositionLine(line string) (*Position, error) {
	if line == "" {
		return nil, nil
	}

	parts := strings.Split(line, ",")

	i, err := parseTime(parts[1])
	if err != nil {
		return nil, err
	}

	x, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, err
	}
	y, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, err
	}
	z, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, err
	}

	result := Position{
		ID:    parts[0],
		Epoch: *i,
		X:     x,
		Y:     y,
		Z:     z,
	}
	return &result, nil
}

// Try a couple of different time formats to make it easier to deal with other languages files.
func parseTime(in string) (*time.Time, error) {
	i, err := time.Parse(time.RFC3339, in)
	if err != nil {
		i, err = time.Parse("2006-01-02T15:04:05.999", in) // default python date format. (milli seconds and no tz)
		if err != nil {
			return nil, err
		}
	}
	return &i, nil
}
