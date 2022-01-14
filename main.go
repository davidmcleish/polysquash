package main

import (
	"fmt"

	"github.com/davidmcleish/polysquash/huffman"
	"github.com/davidmcleish/polysquash/poly"
)

func main() {
	points := []poly.Point{
		{X: 0, Y: 0}, {X: 10, Y: 6}, {X: 7, Y: 11}, {X: -3, Y: 5},
	}
	var tokens []string
	for i := range points {
		g := poly.Gradient(poly.Sub(points[i], points[(i+1)%len(points)]))
		fmt.Println(g)
		tokens = append(tokens, g...)
	}
	table := huffman.CreateHuffman(tokens)
	fmt.Println(table)
	// TODO: origin, scale
}
