package huffman

import "fmt"

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
		exp--
	}
	for val < max {
		val *= 2
		exp++
		if exp > 16 {
			break
		}
	}
	sig := int(val)
	return []string{fmt.Sprint(sign, sig), fmt.Sprint(exp)}
}
