package canon_d

import (
	"fmt"
	"testing"
)

func TestMakeLine(t *testing.T) {
	type TestSet struct {
		in  string
		out string
	}

	sets := []TestSet{
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
	type TestSet struct {
		content string
		before  string
		part    string
		after   string
		found   bool
	}

	sets := []TestSet{
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
