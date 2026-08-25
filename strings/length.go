package strings

func CountRunes(s string) int {
	var count int
	for _, _ = range s {
		count++
	}
	return count
}

func LengthInRange(s string, min, max int) bool {
	return (CountRunes(s) >= min) && (CountRunes(s) <= max)
}
