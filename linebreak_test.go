package linebreak_test

import (
	"fmt"
	"testing"

	"github.com/dgryski/go-linebreak"
)

func ExampleWrap() {
	text := "a b c d e f g h i j k l m n o p qqqqqqqqq"
	width := 9
	textWrapped := linebreak.Wrap(text, width, width)

	fmt.Println(textWrapped)

	// Output:
	// a b c d
	// e f g h
	// i j k l
	// m n o p
	// qqqqqqqqq
}

func TestWrap(t *testing.T) {
	cases := []struct {
		input      string
		width      int
		wantOutput string
	}{
		{
			input: "a b c d e f g h i j k l m n o p qqqqqqqqq",
			width: 9,
			wantOutput: `a b c d
e f g h
i j k l
m n o p
qqqqqqqqq`,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("wrap %q width %d", c.input, c.width), func(t *testing.T) {
			output := linebreak.Wrap(c.input, c.width, c.width)
			if output == c.wantOutput {
				t.Logf("got:\n%s", output)
			} else {
				t.Errorf("got:\n%s\nwant:\n%s\n", output, c.wantOutput)
			}
		})
	}
}
