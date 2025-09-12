package l10n

import (
	stdstrings "strings"
	"unicode"
	"unicode/utf8"

	"github.com/anton2920/gofa/trace"
)

type Language int32

type Localizations map[string][LanguageCount]string

const (
	LanguageEnglish = Language(iota)
	LanguageRussian
	LanguageFrench
	LanguageCount
)

var Language2String = [...]string{
	LanguageEnglish: "English",
	LanguageRussian: "Русский",
	LanguageFrench:  "Français",
}

var Language2HTMLLang = [...]string{
	LanguageEnglish: "en",
	LanguageRussian: "ru",
	LanguageFrench:  "fr",
}

var localizations = make(Localizations)

func (l Language) L(s string) string {
	t := trace.Begin("")

	if l == LanguageEnglish {
		trace.End(t)
		return s
	}

	ls := localizations[s]
	if (len(ls) == 0) || (len(ls[l]) == 0) {
		/*
			switch s {
			default:
				log.Errorf("Not localized %q", s)
			case "↑", "↓", "^|", "|v", "-", "Command":
			}
		*/
		trace.End(t)
		return s
	}

	trace.End(t)
	return ls[l]
}

func Add(ls Localizations) {
	for k, v := range ls {
		k := stdstrings.ToLower(k)

		lowers := v
		for i := 0; i < len(lowers); i++ {
			lowers[i] = stdstrings.ToLower(lowers[i])
		}
		localizations[k] = lowers

		capitals := v
		for i := 0; i < len(capitals); i++ {
			if len(capitals[i]) > 0 {
				r, size := utf8.DecodeRuneInString(capitals[i])
				capitals[i] = string(unicode.ToUpper(r)) + capitals[i][size:]
			}
		}
		localizations[stdstrings.ToTitle(k[:1])+k[1:]] = capitals
	}
}
