package strings

import "testing"

func TestTrim(t *testing.T) {
	samples := [...]struct {
		Input    string
		Expected string
	}{
		{"   ", ""},
		{"  abc   ", "bc"},
		{"  a   b   ", "b"},
	}
	for _, sample := range samples {
		actual := Trim(sample.Input, " a")
		if actual != sample.Expected {
			t.Errorf("Expected %q, got %q", sample.Expected, actual)
		}
	}
}

func TestTrimSpace(t *testing.T) {
	samples := [...]struct {
		Input    string
		Expected string
	}{
		{"   ", ""},
		{"  abc   ", "abc"},
		{"  a   b   ", "a   b"},
	}
	for _, sample := range samples {
		actual := TrimSpace(sample.Input)
		if actual != sample.Expected {
			t.Errorf("Expected %q, got %q", sample.Expected, actual)
		}
	}
}
