package html

type Theme struct {
	/* HTML tags. */
	A        Attributes
	Body     Attributes
	Button   Attributes
	Checkbox Attributes
	Div      Attributes
	Form     Attributes
	H1       Attributes
	H2       Attributes
	H3       Attributes
	H4       Attributes
	H5       Attributes
	H6       Attributes
	Img      Attributes
	Input    Attributes
	LI       Attributes
	Label    Attributes
	OL       Attributes
	P        Attributes
	Select   Attributes
	Span     Attributes
	Textarea Attributes
	UL       Attributes

	/* Default CSS and JS. */
	HeadLink   Attributes
	HeadScript Attributes

	/* Custom components. */
	Error Attributes

	PageSelector                Attributes
	PageSelectorButton          Attributes
	PageSelectorButtonActive    Attributes
	PageSelectorButtonContainer Attributes
}

var nulTheme Theme
