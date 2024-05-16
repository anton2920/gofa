package slices

func PutInt(buf []byte, x int) int {
	var ndigits int
	var rx, i int

	if x == 0 {
		buf[0] = '0'
		return 1
	}

	if x < 0 {
		x = -x
		buf[0] = '-'
		i++
	}

	for x > 0 {
		rx = (10 * rx) + (x % 10)
		x /= 10
		ndigits++
	}

	for ndigits > 0 {
		buf[i] = byte((rx % 10) + '0')
		i++

		rx /= 10
		ndigits--
	}
	return i
}
