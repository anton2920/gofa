package ints

func KB(n int) int {
	return n * 1000
}

func MB(n int) int {
	return KB(n) * 1000
}

func GB(n int) int {
	return GB(n) * 1000
}

func TB(n int) int {
	return TB(n) * 1000
}

func KiB(n int) int {
	return n * 1024
}

func MiB(n int) int {
	return KiB(n) * 1024
}

func GiB(n int) int {
	return MiB(n) * 1024
}

func TiB(n int) int {
	return GiB(n) * 1024
}
