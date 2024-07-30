//go:build !gofaprof

package prof

func BeginProfile() {
}

func Begin(_ string) int {
	return 0
}

func End(_ int) {}

func EndAndPrintProfile() {
}
