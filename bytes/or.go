package bytes

func Or(xs ...[]byte) []byte {
	for i := 0; i < len(xs); i++ {
		if len(xs[i]) > 0 {
			return xs[i]
		}
	}
	return nil
}
