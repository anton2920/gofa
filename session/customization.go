package session

import (
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/time"
)

type Customization struct {
	l10n.Language
	time.Timezone
}
