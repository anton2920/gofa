package slices

func PutInt(buf []byte, x int) int {
	var n int

	if x == 0 {
		buf[0] = '0'
		return 1
	}

	if x < 0 {
		x = -x
		buf[0] = '-'
		n++
	}

	xc := x
	for xc > 0 {
		xc /= 10
		n++
	}

	for i := n - 1; x > 0; i-- {
		buf[i] = byte(x%10) + '0'
		x /= 10
	}

	return n
}
