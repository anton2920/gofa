//go:build !gofaprof

package prof

func BeginProfile() {
}

func Begin(label string) Block {
	return Block{}
}

func End(b Block) {}

func EndAndPrintProfile() {
}
