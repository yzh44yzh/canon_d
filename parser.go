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
			line := MakeLine(l)
			lines = append(lines, line)
		}
	}

	return Card{
		Header: header,
		Lines:  lines,
	}
}

func MakeLine(content string) Line {
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
	idxB := reBold.FindStringIndex(content)
	idxI := reItalic.FindStringIndex(content)
	idxC := reCode.FindStringIndex(content)

	var minIdx []int = nil
	cut := 0

	if idxB != nil {
		minIdx = idxB
		cut = 2
	}

	if idxI != nil {
		if minIdx == nil {
			minIdx = idxI
			cut = 1
		} else {
			if idxI[0] < minIdx[0] {
				minIdx = idxI
				cut = 1
			}
		}
	}

	if idxC != nil {
		if minIdx == nil {
			minIdx = idxC
			cut = 1
		} else {
			if idxC[0] < minIdx[0] {
				minIdx = idxC
				cut = 1
			}
		}
	}

	if minIdx == nil {
		return "", "", "", false
	}

	f := minIdx[0]
	t := minIdx[1]

	return content[:f], content[f+cut : t-cut], content[t:], true
}
