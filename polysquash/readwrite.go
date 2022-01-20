package polysquash

import "github.com/dgryski/go-bitstream"

type Encoder interface {
	Encode(w *bitstream.BitWriter, pts []Point) error
	Decode(r *bitstream.BitReader) ([]Point, error)
}
