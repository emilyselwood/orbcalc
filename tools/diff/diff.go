package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/emilyselwood/orbcalc/orbcore"
)

func main() {
	inA := flag.String("a", "", "first file to compare")
	inB := flag.String("b", "", "second file to compare")

	flag.Parse()

	if *inA == "" || *inB == "" {
		log.Fatal("Need to have two input files A and B")
	}

	aRows, err := orbcore.ReadPositionFile(*inA)
	if err != nil {
		log.Fatal(err)
	}

	bRows, err := orbcore.ReadPositionFile(*inB)
	if err != nil {
		log.Fatal(err)
	}

	if len(aRows) != len(bRows) {
		log.Fatalf("Files contain different numbers of records a: %v b: %v", len(aRows), len(bRows))
	}

	sum := 0.0
	count := 0.0
	for i, a := range aRows {
		d := compareRows(a, bRows[i])
		fmt.Printf("%v: %v\n", count, d/3)
		count++
		sum += d
	}

	log.Printf("difference sum: %v count: %v ave: %v\n", sum, count, sum/(count*3))
}

func compareRows(a *orbcore.Position, b *orbcore.Position) float64 {
	return math.Abs(a.X-b.X) + math.Abs(a.Y-b.Y) + math.Abs(a.Z-b.Z)
}
