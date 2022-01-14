package main

import (
	"fmt"

	"github.com/davidmcleish/polysquash/huffman"
	"github.com/davidmcleish/polysquash/poly"
)

func main() {
	tokens := []string{"a", "a", "a", "a", "a", "a", "a", "b", "b", "b", "b", "c", "d", "e", "f"}
	table := huffman.CreateHuffman(tokens)
	fmt.Println(table)

	points := []poly.Point{
		{X: 0, Y: 0}, {X: 10, Y: 6}, {X: 7, Y: 11}, {X: -3, Y: 5},
	}
	for i := range points {
		fmt.Println(poly.Gradient(poly.Sub(points[i], points[(i+1)%len(points)])))
	}
}
