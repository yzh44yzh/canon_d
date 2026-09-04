package canon_d

import (
	"fmt"
	"testing"
)

func Test_makeCardHeader(t *testing.T) {
	sets := []struct {
		in    string
		out   string
		level int
	}{
		{in: "# aaa ", out: "aaa", level: 1},
		{in: "#    aaa bb ", out: "aaa bb", level: 1},
		{in: "##  bbbb ", out: "bbbb", level: 2},
		{in: "### ccc", out: "ccc", level: 3},
		{in: "##### Some Header", out: "Some Header", level: 5},
	}

	for i, set := range sets {
		name := fmt.Sprintf("set_%d", i)
		t.Run(name, func(t *testing.T) {
			header := makeCardHeader(set.in)
			if header.Content != set.out {
				t.Errorf("in: '%v', expect out: '%v', real out: '%v'", set.in, set.out, header.Content)
			}
			if header.Level != set.level {
				t.Errorf("in: '%v', expect level: '%v', real level: '%v'", set.in, set.level, header.Level)
			}
		})
	}
}

func TestCardHeaderShow(t *testing.T) {
	h1 := CardHeader{
		Level:   1,
		Content: "aaa",
	}
	h2 := CardHeader{
		Level:   2,
		Content: "bbb",
		Parent:  &h1,
	}
	h3 := CardHeader{
		Level:   3,
		Content: "ccc",
		Parent:  &h2,
	}

	s1 := h1.Show()
	if s1 != "aaa" {
		t.Errorf("invalid result for h1.Show '%v'", s1)
	}

	s2 := h2.Show()
	if s2 != "aaa / bbb" {
		t.Errorf("invalid result for h2.Show '%v'", s2)
	}

	s3 := h3.Show()
	if s3 != "aaa / bbb / ccc" {
		t.Errorf("invalid result for h3.Show '%v'", s3)
	}
}

func TestMakeLine(t *testing.T) {
	sets := []struct {
		in  string
		out string
	}{
		{in: "a a a", out: "a a a"},
		{in: "a `b` c", out: "a ??? c"},
		{in: "a `b` c `d` e", out: "a ??? c ??? e"},
		{in: "`b` c `d` e", out: "??? c ??? e"},
		{in: "ha `he` hoo", out: "ha ??? hoo"},
		{in: "a `b` c `d`", out: "a ??? c ???"},
		{in: "`b`", out: "???"},
		{in: "aa `bb` cc", out: "aa ??? cc"},
		{in: "aa **bb** cc", out: "aa ??? cc"},
		{in: "aaa **bbb** ccc **ddd** eee", out: "aaa ??? ccc ??? eee"},
		{in: "**b** c `d` e", out: "??? c ??? e"},
		{in: "ha **he** hoo", out: "ha ??? hoo"},
		{in: "a `b` c **d**", out: "a ??? c ???"},
		{in: "**b**", out: "???"},
		{in: "a *b* c *d*", out: "a *b* c *d*"},
		{in: "**b** c _d_ e", out: "??? c ??? e"},
	}

	for i, set := range sets {
		name := fmt.Sprintf("set_%d", i)
		t.Run(name, func(t *testing.T) {
			line := makeLine(set.in)
			if line.Show() != set.out {
				t.Errorf("in: '%v', expect out: '%v', real out: '%v', line: '%v'", set.in, set.out, line.Show(), line)
			}
		})
	}
}

func TestFindPart(t *testing.T) {
	sets := []struct {
		content string
		before  string
		part    string
		after   string
		found   bool
	}{
		{
			content: "Семейство движков `MergeTree` позволяют",
			before:  "Семейство движков ",
			part:    "MergeTree",
			after:   " позволяют",
			found:   true,
		},
		{
			content: "Семейство движков MergeTree позволяют",
			before:  "",
			part:    "",
			after:   "",
			found:   false,
		},
		{
			content: "**aaa** bbb",
			before:  "",
			part:    "aaa",
			after:   " bbb",
			found:   true,
		},
		{
			content: "aaa **bbb**",
			before:  "aaa ",
			part:    "bbb",
			after:   "",
			found:   true,
		},
		{
			content: "`aaa` bbb ccc",
			before:  "",
			part:    "aaa",
			after:   " bbb ccc",
			found:   true,
		},
		{
			content: "aaa _bbb_ ccc",
			before:  "aaa ",
			part:    "bbb",
			after:   " ccc",
			found:   true,
		},
		{
			content: "`aaa` _bbb_ **ccc**",
			before:  "",
			part:    "aaa",
			after:   " _bbb_ **ccc**",
			found:   true,
		},
		{
			content: "aaa **bbb** `ccc`",
			before:  "aaa ",
			part:    "bbb",
			after:   " `ccc`",
			found:   true,
		},
	}

	for i, set := range sets {
		name := fmt.Sprintf("set_%d", i)
		t.Run(name, func(t *testing.T) {
			before, part, after, found := findPart(set.content)
			if before != set.before {
				t.Errorf("%d expect before: '%v', real before: '%v'", i, set.before, before)
			}
			if part != set.part {
				t.Errorf("%d expect part: '%v', real part: '%v'", i, set.part, part)
			}
			if after != set.after {
				t.Errorf("%d expect after: '%v', real after: '%v'", i, set.after, after)
			}
			if found != set.found {
				t.Errorf("%d expect found: '%v', real found: '%v'", i, set.found, found)
			}
		})
	}
}
