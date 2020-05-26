package parser

import (
	"strings"
	"testing"

	"lukechampine.com/go-urbit/hoon/ast"
	"lukechampine.com/go-urbit/hoon/scanner"
)

func TestRoundTrip(t *testing.T) {
	var tests = []struct {
		prog string
		exp  string
	}{
		{
			prog: `[1 2 3]`,
			exp:  `[1 [2 3]]`,
		},
		{
			prog: `=/  n  1
                   n`,
			exp: `=/(n 1 n)`,
		},
		{
			prog: `=/  a  2
                   =/  b  7
                   ?:  =(a +(b))
                     a
                   b`,
			exp: `=/(a 2 =/(b 7 ?:(.=(a .+(b)) a b)))`,
		},
		{
			prog: `=/  n  1
                   [. .]:n`,
			exp: `=/(n 1 =<([. .] n))`,
		},
		{
			prog: `=/  n  0
                   |-
                   ?:  =(n 5)
                     n
                   $(n +(n))`,
			exp: `=/(n 0 |-(?:(.=(n 5) n %=($ n .+(n)))))`,
		},
		{
			prog: `|=  a=@
                   =/  b  2
                   =/  f  |=(@ 7)
                   (f(a 2, b 3))`,
			exp: `|=(a=@ =/(b 2 =/(f |=(@ 7) (%=(f a 2, b 3)))))`,
		},
		{
			prog: `|=  n=@
                   =/  acc=@  1
                   |-
                   ?:  =(n 0)  acc
                   %=  $
                     n  (dec n)
                     acc  (mul acc n)
                   ==`,
			exp: `|=(n=@ =/(acc=@ 1 |-(?:(.=(n 0) acc %=($ n (dec n), acc (mul acc n))))))`,
		},
	}
	for _, test := range tests {
		var sb strings.Builder
		ast.Print(&sb, New(scanner.New([]byte(test.prog))).Parse())
		if got := sb.String(); got != test.exp {
			t.Fatalf("bad parse:\nexp: %q\ngot: %q", test.exp, got)
		}
	}
}
