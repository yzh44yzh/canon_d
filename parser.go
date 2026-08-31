package canon_d

import (
	"regexp"
	"strings"
)

// **bold text**, but __bold text__ is not supported
var reBold *regexp.Regexp = regexp.MustCompile(`\*\*(\S.*?)\*\*`)

// _italic text_, but *italic text* is not supported
var reItalic *regexp.Regexp = regexp.MustCompile(`\_(\S.*?)\_`)

// `inline code`
var reCode *regexp.Regexp = regexp.MustCompile("`(\\S.*?)`")

func MakeCard(content string) Card {
	rawLines := strings.Split(content, "\n")
	header := ""
	lines := []Line{}

	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		if header == "" {
			header = strings.Trim(l, "#")
			header = strings.TrimSpace(header)
		} else {
			line := makeLine(l)
			lines = append(lines, line)
		}
	}

	return Card{
		Header: header,
		Lines:  lines,
	}
}

func makeLine(content string) Line {
	lineParts := []LinePart{}
	rest := content

	for {
		before, part, after, found := findPart(rest)
		if found {
			if before != "" {
				beforePart := LinePart{
					Content: before,
					Visible: true,
				}
				lineParts = append(lineParts, beforePart)
			}

			currPart := LinePart{
				Content: part,
				Visible: false,
			}
			lineParts = append(lineParts, currPart)

			rest = after
		} else {
			if rest != "" {
				restPart := LinePart{
					Content: rest,
					Visible: true,
				}
				lineParts = append(lineParts, restPart)
			}
			break
		}
	}

	return Line{
		Original: content,
		Parts:    lineParts,
	}
}

func findPart(content string) (before, part, after string, found bool) {
	var minIdx []int = nil
	cut := 0

	idxB := reBold.FindStringIndex(content)
	if idxB != nil {
		minIdx = idxB
		cut = 2
	}

	idxI := reItalic.FindStringIndex(content)
	switch {
	case idxI != nil && minIdx == nil:
		minIdx = idxI
		cut = 1
	case idxI != nil && idxI[0] < minIdx[0]:
		minIdx = idxI
		cut = 1
	}

	idxC := reCode.FindStringIndex(content)
	switch {
	case idxC != nil && minIdx == nil:
		minIdx = idxC
		cut = 1
	case idxC != nil && idxC[0] < minIdx[0]:
		minIdx = idxC
		cut = 1
	}

	if minIdx == nil {
		return "", "", "", false
	}

	f := minIdx[0]
	t := minIdx[1]

	return content[:f], content[f+cut : t-cut], content[t:], true
}
