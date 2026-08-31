package canon_d

import (
	"testing"
)

func Test_MakeLine(t *testing.T) {
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

	for _, set := range sets {
		line := MakeLine(set.in)
		if line.Show() != set.out {
			t.Errorf("in: '%v', expect out: '%v', real out: '%v', line: '%v'", set.in, set.out, line.Show(), line)
		}
	}
}
