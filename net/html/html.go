package html

const Header = `<!DOCTYPE html>`

var (
	Quot = "&#34;" // shorter than "&quot;"
	Apos = "&#39;" // shorter than "&apos;" and apos was not in HTML until HTML5
	Amp  = "&amp;"
	Lt   = "&lt;"
	Gt   = "&gt;"
	Null = "\uFFFD"
)
