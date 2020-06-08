package utils

import "strconv"

func DivideAsPercent(a int64, b int64) string {
	if a == 0 || b == 0 {
		return "0 %"
	}
	return strconv.FormatFloat(float64(a)/float64(b)*100, 'f', 4, 32) + " %"
}
