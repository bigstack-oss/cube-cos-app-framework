package base

var (
	Welcome          bool   = true
	Header           bool   = true
	Format           string = "text"
	NameOnly         bool   = false
	SupportedFormats        = map[string]struct{}{
		"text": {},
		"json": {},
	}
)
