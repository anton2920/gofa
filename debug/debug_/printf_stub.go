//go:build !gofadebug
// +build !gofadebug

package debug_

func Printf(_ string, _ ...interface{}) (int, error) {
	return 0, nil
}
