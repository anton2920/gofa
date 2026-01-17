package strings

import "strconv"

func ToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
