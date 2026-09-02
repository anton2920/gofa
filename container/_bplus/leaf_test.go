package bplus

import "testing"

func BenchmarkLeafInsertKeyValueAt(b *testing.B) {
	var page Page
	page.Init(PageTypeLeaf)

	leaf := page.Leaf()
	key := Uint16ToBytes(0)
	value := Uint16ToBytes(0)

	b.Run("Prepend", func(b *testing.B) {
		leaf.InsertKeyValueAt(key, value, 0)
		for i := 0; i < b.N; i++ {
			if leaf.OverflowAfterInsertKeyValue(len(key), len(value)) {
				leaf.Reset()
			}
			leaf.InsertKeyValueAt(key, value, 0)
		}
	})
	b.Run("Append", func(b *testing.B) {
		leaf.InsertKeyValueAt(key, value, 0)
		for i := 0; i < b.N; i++ {
			if leaf.OverflowAfterInsertKeyValue(len(key), len(value)) {
				leaf.Reset()
			}
			leaf.InsertKeyValueAt(key, value, int(leaf.N))
		}
	})
}
