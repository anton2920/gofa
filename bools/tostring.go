package bools

var strings = [2]string{"false", "true"}

//go:nosplit
func ToString(b bool) string {
	return strings[ToInt(b)]
}
