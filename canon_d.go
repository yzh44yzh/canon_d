package canon_d

import (
	"strings"
)

type Deck struct {
	Header string
	Cards  []Card
}

type Card struct {
	Header string
	Lines  []Line
}

type Line struct {
	Original string
	Parts    []LinePart
}

type LinePart struct {
	Visible bool
	Content string
}

func (l Line) Show() string {
	var res strings.Builder
	for _, part := range l.Parts {
		if part.Visible {
			res.WriteString(part.Content)
		} else {
			res.WriteString("???")
		}
	}
	return res.String()
}
