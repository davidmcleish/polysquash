package main

import (
	"fmt"

	"github.com/davidmcleish/polysquash/huffman"
	"github.com/davidmcleish/polysquash/poly"
)

func main() {
	points := []poly.Point{
		{X: 2000, Y: 3000}, {X: 2010, Y: 3006}, {X: 2007, Y: 3001}, {X: 1997, Y: 2995},
	}
	tokens := poly.Quantise(points)
	fmt.Println(tokens)

	table := huffman.CreateHuffman(tokens)
	fmt.Println(table)

	parsed, err := poly.Parse(tokens)
	fmt.Println(parsed, err)
}
