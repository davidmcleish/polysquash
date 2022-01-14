package huffman

import (
	"fmt"
	"math"
)

type Code string

type tree struct {
	freq        int
	leaf        string
	left, right *tree
}

func (t *tree) String() string {
	return fmt.Sprintf("%+v", *t)
}

func CreateHuffman(tokens []string) map[string]Code {
	if len(tokens) == 0 {
		return map[string]Code{}
	}
	// TODO: index
	var forest []*tree
	for _, tok := range tokens {
		found := false
		for i, t := range forest {
			if t.leaf == tok {
				forest[i].freq++
				found = true
				break
			}
		}
		if !found {
			forest = append(forest, &tree{1, tok, nil, nil})
		}
	}
	// TODO: heap
	for len(forest) > 1 {
		fMin1 := math.MaxInt32
		fMin2 := math.MaxInt32
		iMin1 := -1
		iMin2 := -1
		for i, t := range forest {
			if t.freq < fMin2 {
				fMin2 = t.freq
				iMin2 = i
			}
			if fMin2 < fMin1 {
				fMin1, fMin2 = fMin2, fMin1
				iMin1, iMin2 = iMin2, iMin1
			}
		}
		if iMin1 > iMin2 {
			iMin1, iMin2 = iMin2, iMin1
		}
		forest[iMin1] = &tree{
			freq:  fMin1 + fMin2,
			left:  forest[iMin1],
			right: forest[iMin2],
		}
		forest[iMin2] = forest[len(forest)-1]
		forest = forest[:len(forest)-1]
	}
	table := make(map[string]Code)
	forest[0].buildTable(table, "")
	return table
}

func (t *tree) buildTable(table map[string]Code, prefix Code) {
	if t.left == nil && t.right == nil {
		table[t.leaf] = prefix
	} else {
		t.left.buildTable(table, prefix+"0")
		t.right.buildTable(table, prefix+"1")
	}
}
