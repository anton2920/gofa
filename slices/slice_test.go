package slices

import (
	"crypto/rand"
	"strconv"
	"testing"

	"github.com/anton2920/gofa/util"
)

func TestPutInt(t *testing.T) {
	buffer := make([]byte, 32)

	for i := 0; i < 100_000; i++ {
		expected := rand.Int()
		n := PutInt(buffer, expected)
		actual, err := strconv.Atoi(util.Slice2String(buffer[:n]))
		if err != nil {
			t.Errorf("Failed to decode resulting integer; expected %d, got %s", expected, buffer[:n])
		}
		if expected != actual {
			t.Errorf("Expected to decode %d, got %d", expected, actual)
		}
	}
}
