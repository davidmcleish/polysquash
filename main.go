package main

import (
	"fmt"

	"github.com/davidmcleish/polysquash/huffman"
)

func main() {
	tokens := []huffman.Token{"a", "a", "a", "a", "a", "a", "a", "b", "b", "b", "b", "c", "d", "e", "f"}
	table := huffman.CreateHuffman(tokens)
	fmt.Println(table)
}
