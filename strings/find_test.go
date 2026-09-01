package strings

import "testing"

func TestFindSubstring(t *testing.T) {
	samples := [...]struct {
		Haystack string
		Needle   string
		Expected int
	}{
		{"hh := html.New(w, r); h := &hh}}", "}}", 30},
	}
	for _, sample := range samples {
		actual := FindSubstring(sample.Haystack, sample.Needle)
		if actual != sample.Expected {
			t.Errorf("for haystack %q and needle %q expected %d, got %d", sample.Haystack, sample.Needle, sample.Expected, actual)
		}
	}
}
