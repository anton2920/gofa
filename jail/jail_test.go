package jail

import "testing"

func BenchmarkCreateRemoveJail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jail, err := NewJail()
		if err != nil {
			b.Error("Failed to create new jail: ", err)
		}
		if err := RemoveJail(jail); err != nil {
			b.Error("Failed to remove jail: ", err)
		}
	}
}
