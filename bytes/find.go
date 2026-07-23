package bytes

func FindByte(xs []byte, x byte) int {
	for i := 0; i < len(xs); i++ {
		if xs[i] == x {
			return i
		}
	}
	return -1
}

func FindByteReverse(xs []byte, x byte) int {
	for i := len(xs) - 1; i >= 0; i-- {
		if xs[i] == x {
			return i
		}
	}
	return -1
}
