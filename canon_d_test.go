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
		{in: "ha `he` hoo", out: "he he hoo"},
		{in: "a `b` c `d`", out: "a ??? c ???"},
		{in: "`b`", out: "???"},
	}

	for _, set := range sets {
		line := MakeLine(set.in)
		if line.Show() != set.out {
			t.Errorf("incorrect Line '%v' for input '%v'", line, set.in)
		}
	}
}
