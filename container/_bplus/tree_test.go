package bplus

import (
	"bytes"
	"crypto/rand"
	"testing"
)

const N = 10000

func testTreeGet(t *testing.T, g Generator, pager Pager) {
	t.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		t.Fatalf("Failed to create new tree: %v", err)
	}

	m := make(map[int]int)
	for i := 0; i < N; i++ {
		k := g.Generate()
		v := g.Generate()

		m[k] = v
		tree.Set(int2Slice(k), int2Slice(v))
	}

	for k, v := range m {
		got, err := tree.Get(int2Slice(k))
		if err != nil {
			t.Errorf("Error on 'Get': %v", err)
		} else if slice2Int(got) != v {
			t.Errorf("Expected value %v, got %v", v, slice2Int(got))
		}
	}
}

func testTreeDel(t *testing.T, g Generator, pager Pager) {
	t.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		t.Fatalf("Failed to create new tree: %v", err)
	}

	m := make(map[int]struct{})
	for i := 0; i < N; i++ {
		k := g.Generate()

		m[k] = struct{}{}
		if err := tree.Set(int2Slice(k), ZeroValue); err != nil {
			t.Fatalf("Error on 'Set': %v", err)
		}
	}

	for k := range m {
		if err := tree.Del(int2Slice(k)); err != nil {
			t.Errorf("Error on 'Del': %v", err)
		}

		ok, err := tree.Has(int2Slice(k))
		if err != nil {
			t.Fatalf("Error on 'Has': %v", err)
		} else if ok {
			t.Errorf("Expected key %v to be removed, but it's still present", k)
		}
	}
}

func testTreeHas(t *testing.T, g Generator, pager Pager) {
	t.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		t.Fatalf("Failed to create new tree: %v", err)
	}

	m := make(map[int]struct{})
	for i := 0; i < N; i++ {
		k := g.Generate()

		m[k] = struct{}{}
		tree.Set(int2Slice(k), ZeroValue)
	}

	for k := range m {
		ok, err := tree.Has(int2Slice(k))
		if err != nil {
			t.Errorf("Error on 'Has': %v", err)
		} else if !ok {
			t.Errorf("Expected to find key %v, found nothing", k)
		}
	}
}

func testTreeSet(t *testing.T, g Generator, pager Pager) {
	t.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		t.Fatalf("Failed to create new tree: %v", err)
	}

	for i := 0; i < N; i++ {
		k := g.Generate()
		v := g.Generate()

		if err := tree.Set(int2Slice(k), int2Slice(v)); err != nil {
			t.Errorf("Error on 'Set': %v", err)
		}

		ok, err := tree.Has(int2Slice(k))
		if err != nil {
			t.Fatalf("Error on 'Has': %v", err)
		} else if !ok {
			t.Errorf("Expected to find key %v, found nothing", k)
		}

		got, err := tree.Get(int2Slice(k))
		if err != nil {
			t.Fatalf("Error on 'Get': %v", err)
		} else if slice2Int(got) != v {
			t.Errorf("Expected value %v, got %v", v, slice2Int(got))
		}
	}
}

func testTreeSetLarge(t *testing.T, g Generator, pager Pager) {
	t.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		t.Fatalf("Failed to create new tree: %v", err)
	}
	value := make([]byte, PageSize)

	for i := 0; i < N; i++ {
		k := g.Generate()
		if _, err := rand.Read(value); err != nil {
			t.Fatalf("Failed to generate random value: %v", err)
		}

		if err := tree.Set(int2Slice(k), value); err != nil {
			t.Errorf("Error on 'Set': %v", err)
		}

		ok, err := tree.Has(int2Slice(k))
		if err != nil {
			t.Fatalf("Error on 'Has': %v", err)
		} else if !ok {
			t.Errorf("Expected to find key %v, found nothing", k)
		}

		got, err := tree.Get(int2Slice(k))
		if err != nil {
			t.Fatalf("Error on 'Get': %v", err)
		} else if bytes.Compare(got, value) != 0 {
			t.Errorf("Expected value %v, got %v", value, got)
		}
	}
}

func TestTree(t *testing.T) {
	ops := [...]struct {
		Name string
		Func func(*testing.T, Generator, Pager)
	}{
		{"Get", testTreeGet},
		// 	{"Del", testTreeDel},
		{"Has", testTreeHas},
		{"Set", testTreeSet},
		{"SetLarge", testTreeSetLarge},
	}

	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}

	for _, op := range ops {
		t.Run(op.Name, func(t *testing.T) {
			for _, generator := range generators {
				generator.Reset()
				t.Run(generator.String(), func(t *testing.T) {
					t.Parallel()
					t.Run("MemoryPager", func(t *testing.T) {
						op.Func(t, generator, new(MemoryPager))
					})
					/*
						t.Run("FilePager", func(t *testing.T) {
							t.Skip()
							filePager, err := FilePagerNew(generator.String() + "_test.tree")
							if err != nil {
								t.Fatalf("Failed to create new file pager: %v", err)
							}
							defer filePager.Close()
							op.Func(t, generator, filePager)
						})
					*/
				})
			}
		})
	}
}

func benchmarkTreeGet(b *testing.B, g Generator, pager Pager) {
	b.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		b.Fatalf("Failed to create new tree: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_ = tree.Set(int2Slice(g.Generate()), ZeroValue)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tree.Get(int2Slice(g.Generate()))
	}
}

func benchmarkTreeDel(b *testing.B, g Generator, pager Pager) {
	b.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		b.Fatalf("Failed to create new tree: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_ = tree.Set(int2Slice(g.Generate()), ZeroValue)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tree.Del(int2Slice(g.Generate()))
	}
}

func benchmarkTreeSet(b *testing.B, g Generator, pager Pager) {
	b.Helper()

	tree, err := GetTreeAt(pager, -1)
	if err != nil {
		b.Fatalf("Failed to create new tree: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_ = tree.Set(int2Slice(g.Generate()), ZeroValue)
	}
}

func BenchmarkTree(b *testing.B) {
	ops := [...]struct {
		Name string
		Func func(*testing.B, Generator, Pager)
	}{
		{"Get", benchmarkTreeGet},
		//	{"Del", benchmarkTreeDel},
		{"Set", benchmarkTreeSet},
	}

	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}

	for _, op := range ops {
		b.Run(op.Name, func(b *testing.B) {
			for _, generator := range generators {
				generator.Reset()
				b.Run(generator.String(), func(b *testing.B) {
					b.Run("MemoryPager", func(b *testing.B) {
						op.Func(b, generator, new(MemoryPager))
					})
					/*
						b.Run("FilePager", func(b *testing.B) {
							filePager, err := FilePagerNew(generator.String() + "_test.tree")
							if err != nil {
								b.Fatalf("Failed to create new file pager: %v", err)
							}
							defer filePager.Close()
							op.Func(b, generator, filePager)
						})
					*/
				})
			}
		})
	}
}
