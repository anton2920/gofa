package l10n

import (
	stdstrings "strings"

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

		lower := v
		for i := 0; i < len(lower); i++ {
			lower[i] = stdstrings.ToLower(lower[i])
		}
		localizations[k] = lower

		title := v
		for i := 0; i < len(title); i++ {
			title[i] = stdstrings.Title(title[i])
		}
		localizations[stdstrings.Title(k)] = title
	}
}
