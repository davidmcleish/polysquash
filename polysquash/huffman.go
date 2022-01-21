package polysquash

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/dgryski/go-bitstream"
	"github.com/peterstace/simplefeatures/geom"
)

type Code string

type tree struct {
	freq        int
	leaf        int64
	left, right *tree
}

func (t *tree) String() string {
	return fmt.Sprintf("%+v", *t)
}

func CreateHuffman(tokens []int64) *tree {
	if len(tokens) == 0 {
		return nil
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
	return forest[0]
}

func (t *tree) buildTable(table map[int64]Code, prefix Code) {
	if t.left == nil && t.right == nil {
		table[t.leaf] = prefix
	} else {
		t.left.buildTable(table, prefix+"0")
		t.right.buildTable(table, prefix+"1")
	}
}

func (t *tree) readToken(r *bitstream.BitReader) (int64, error) {
	if t.left == nil && t.right == nil {
		return t.leaf, nil
	} else {
		bit, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if bit == bitstream.Zero {
			return t.left.readToken(r)
		} else {
			return t.right.readToken(r)
		}
	}
}

func HuffmanEncode(w *bitstream.BitWriter, tokens []int64) error {
	tree := CreateHuffman(tokens)
	table := make(map[int64]Code)
	tree.buildTable(table, "")

	if err := tree.writeHeader(w); err != nil {
		return err
	}

	// write length
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, int64(len(tokens)))
	for _, b := range buf[:n] {
		if err := w.WriteByte(b); err != nil {
			return err
		}
	}

	for _, tok := range tokens {
		code, ok := table[tok]
		if !ok {
			panic("token not in table")
		}
		for _, c := range code {
			bit := bitstream.Zero
			if c == '1' {
				bit = bitstream.One
			}
			if err := w.WriteBit(bit); err != nil {
				return err
			}
		}
	}
	return nil
}

func HuffmanDecode(r *bitstream.BitReader) ([]int64, error) {
	tree, err := readHeader(r)
	if err != nil {
		return nil, err
	}

	// read length
	length, err := binary.ReadVarint(r)
	if err != nil {
		return nil, err
	}
	tokens := make([]int64, length)

	for i := range tokens {
		var err error
		tokens[i], err = tree.readToken(r)
		if err != nil {
			return nil, err
		}
	}
	return tokens, nil
}

func (t *tree) writeHeader(w *bitstream.BitWriter) error {
	if t.left == nil && t.right == nil {
		if err := w.WriteBit(bitstream.Zero); err != nil {
			return err
		}
		buf := make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(buf, int64(t.leaf))
		for _, b := range buf[:n] {
			if err := w.WriteByte(b); err != nil {
				return err
			}
		}
	} else {
		if err := w.WriteBit(bitstream.One); err != nil {
			return err
		}
		t.left.writeHeader(w)
		t.right.writeHeader(w)
	}
	return nil
}

func readHeader(r *bitstream.BitReader) (*tree, error) {
	bit, err := r.ReadBit()
	if err != nil {
		return nil, err
	}
	if bit == bitstream.Zero {
		leaf, err := binary.ReadVarint(r)
		if err != nil {
			return nil, err
		}
		return &tree{leaf: leaf}, nil
	} else {
		left, err := readHeader(r)
		if err != nil {
			return nil, err
		}
		right, err := readHeader(r)
		if err != nil {
			return nil, err
		}
		return &tree{left: left, right: right}, nil
	}
}

type HuffmanWKT struct{}

func (t HuffmanWKT) String() string { return "HuffWKT" }

func (t HuffmanWKT) Encode(w io.Writer, poly geom.Polygon) error {
	wkt := poly.AsText()
	tokens := make([]int64, len(wkt))
	for i, c := range wkt {
		tokens[i] = int64(c)
	}
	bw := bitstream.NewWriter(w)
	if err := HuffmanEncode(bw, tokens); err != nil {
		return err
	}
	return bw.Flush(bitstream.Zero)
}

func (t HuffmanWKT) Decode(r io.Reader) (*geom.Polygon, error) {
	br := bitstream.NewReader(r)
	tokens, err := HuffmanDecode(br)
	if err != nil {
		return nil, err
	}
	wkt := make([]rune, len(tokens))
	for i, tok := range tokens {
		wkt[i] = rune(tok)
	}

	g, err := geom.UnmarshalWKT(string(wkt))
	if err != nil {
		return nil, err
	}
	poly, ok := g.AsPolygon()
	if !ok {
		return nil, errors.New("not a polygon")
	}
	return &poly, nil
}
