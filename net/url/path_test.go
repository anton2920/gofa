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
	t.Run("/company/1/edit", func(t *testing.T) {
		t.Parallel()

		var id database.ID
		path := Path("/company/1/edit")
		MustMatch(t, &path, "/company...")
		MustMatch(t, &path, "/%d...", &id)
		MustMatch(t, &path, "/edit")
		if id != 1 {
			t.Errorf("expected ID=1, got %d", id)
		}
	})
	t.Run("(json,wire)", func(t *testing.T) {
		t.Parallel()

		const expectedList = "json,wire"
		var actualList string
		path := Path("(" + expectedList + ")")
		MustMatch(t, &path, "(%s)", &actualList)
		if actualList != expectedList {
			t.Errorf("expected list=%q, got %q", expectedList, actualList)
		}
	})
}
