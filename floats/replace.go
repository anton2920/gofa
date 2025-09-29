package floats

func Replace32(r *float32, x float32) {
	if x > 0 {
		*r = x
	}
}

func Replace64(r *float64, x float64) {
	if x > 0 {
		*r = x
	}
}
