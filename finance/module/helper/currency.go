package helper

import (
	"math"
	"strconv"
)

func FormatRupiah(n int64) string {
	rupiah := strconv.FormatInt(int64(math.Abs(float64(n))), 10)

	var rupiahFormatted string
	var idx int
	for i := len(rupiah) - 1; i >= 0; i-- {
		rupiahFormatted = string(rupiah[i]) + rupiahFormatted
		idx++
		if idx%3 == 0 && i != 0 {
			rupiahFormatted = "." + rupiahFormatted
		}
	}

	if n < 0 {
		return "-Rp" + rupiahFormatted
	}

	return "Rp" + rupiahFormatted
}
