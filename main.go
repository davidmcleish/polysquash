package main

import (
	"fmt"

	"github.com/davidmcleish/polysquash/polysquash"
)

func main() {
	points := []polysquash.Point{
		{X: 2000, Y: 3000}, {X: 2010, Y: 3006}, {X: 2007, Y: 3001}, {X: 1997, Y: 2995},
	}
	tokens := polysquash.QuantisePts(points)
	fmt.Println(tokens)

	table := polysquash.CreateHuffman(tokens)
	fmt.Println(table)

	parsed, err := polysquash.Parse(tokens)
	fmt.Println(parsed, err)
}
