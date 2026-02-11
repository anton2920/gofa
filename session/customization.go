package session

import (
	"github.com/anton2920/gofa/l10n"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/time"
)

type ColorScheme int32

type Customization struct {
	l10n.Language
	time.Timezone
	ColorScheme
}

const (
	ColorSchemeAuto = ColorScheme(iota)
	ColorSchemeLight
	ColorSchemeDark
	ColorSchemeCount
)

var ColorScheme2String = [...]string{
	ColorSchemeAuto:  "auto",
	ColorSchemeLight: "light",
	ColorSchemeDark:  "dark",
}

func FillCustomizationFromRequest(vs url.Values, customization *Customization) {
	lang, _ := vs.GetInt("Language")
	customization.Language = l10n.Language(lang)

	tz, _ := vs.GetInt("Timezone")
	customization.Timezone = time.Timezone(tz)

	clr, _ := vs.GetInt("ColorScheme")
	customization.ColorScheme = ColorScheme(clr)
}
