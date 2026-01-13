package url

import (
	"testing"

	"github.com/anton2920/gofa/database"
)

func MustMatch(t *testing.T, path *Path, format string, args ...interface{}) {
	t.Helper()
	if !path.Match(format, args...) {
		t.Errorf("expected match on %q, but didn't get one", format)
	}
}

func TestPathMatch(t *testing.T) {
	name := t.Name() + "/"

	t.Run("/company/1/edit", func(t *testing.T) {
		t.Parallel()

		var id database.ID
		path := Path(t.Name()[len(name):])
		MustMatch(t, &path, "/company...")
		MustMatch(t, &path, "/%d...", &id)
		MustMatch(t, &path, "/edit")
		if id != 1 {
			t.Errorf("expected ID=1, got %d", id)
		}
	})
	t.Run("/company/1/logo.jpg", func(t *testing.T) {
		t.Parallel()

		var id database.ID
		path := Path(t.Name()[len(name):])
		MustMatch(t, &path, "/company...")
		MustMatch(t, &path, "/%d...", &id)
		MustMatch(t, &path, "/logo.jpg")
		if id != 1 {
			t.Errorf("expected ID=1, got %d", id)
		}
	})
	t.Run("(json,wire)", func(t *testing.T) {
		t.Parallel()

		var expectedList = t.Name()[len(name)+1 : len(t.Name())-1]
		var actualList string
		path := Path(t.Name()[len(name):])
		MustMatch(t, &path, "(%s)", &actualList)
		if actualList != expectedList {
			t.Errorf("expected list=%q, got %q", expectedList, actualList)
		}
	})
}
