package fmt

import (
	"os"
	"testing"

	"github.com/tenntenn/golden"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name    string
		fmtcmds map[string]string
		in      string
	}{
		{
			"format inline JavaScript",
			map[string]string{
				"Expr": "deno fmt ${FILE} --ext js",
			},
			"../testdata/deno_fmt.cue",
		},
		{
			"format inline GraphQL",
			map[string]string{
				"Query": "prettier ${FILE} --parser graphql",
			},
			"../testdata/gql_fmt.cue",
		},
		{
			"format inline JavaScript and GraphQL",
			map[string]string{
				"Expr":  "deno fmt ${FILE} --ext js",
				"Query": "prettier ${FILE} --parser graphql",
			},
			"../testdata/deno_gql_fmt.cue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fmtcmds)
			b, err := os.ReadFile(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			formatted, err := c.Format(b)
			if err != nil {
				t.Fatal(err)
			}
			got := string(formatted)
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, "", tt.in, got)
				return
			}
			if diff := golden.Diff(t, "", tt.in, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
