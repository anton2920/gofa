//go:build !trace

package trace

func BeginProfile() {
}

func Start(label string) Block {
	return Block{}
}

func End(b Block) {}

func EndAndPrintProfile() {
}
