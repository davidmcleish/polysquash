package huffman

import (
	"fmt"
	"strconv"

	"github.com/juju/errors"
)

func Quantise(val float64) []string {
	sign := ""
	if val < 0 {
		val = -val
		sign = "-"
	}
	const max = 32768
	exp := 0
	for val > max {
		val /= 2
		exp++
	}
	for val < max {
		val *= 2
		exp--
		if exp > 16 {
			break
		}
	}
	sig := int(val)
	return []string{fmt.Sprint(sign, sig), fmt.Sprint(exp)}
}

func FromQuantised(sigStr, expStr string) (float64, error) {
	sig, err := strconv.Atoi(sigStr)
	if err != nil {
		return 0, errors.Trace(err)
	}
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return 0, errors.Trace(err)
	}
	if exp < 0 {
		return float64(sig) / float64(uint32(1)<<(-exp)), nil
	} else {
		return float64(sig) * float64(uint32(1)<<exp), nil
	}
}
