package slug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlug(t *testing.T) {

	cases := [][]struct {
		in, out string
	}{
		{
			{"foo", "foo"},
			{"foo", "foo-1"},
			{"foo bar", "foo-bar"},
		},
		{
			{"foo", "foo"},
			{"fooCamelCase", "foocamelcase"},
		},
		{
			{"foo", "foo"},
			{"foo", "foo-1"},
			// {"foo 1", "foo-1-1"}, // these are too rare for Jsonnet
			// {"foo 1", "foo-1-2"},
			{"foo", "foo-2"},
		},
		{
			{"heading with a - dash", "heading-with-a---dash"},
			{"heading with an _ underscore", "heading-with-an-_-underscore"},
			{"heading with a period.txt", "heading-with-a-periodtxt"},
			{"exchange.bind_headers(exchange, routing [, bindCallback])", "exchangebind_headersexchange-routing--bindcallback"},
		},
	}

	for _, cs := range cases {
		s := New()
		for _, c := range cs {
			assert.Equal(t, c.out, s.Slug(c.in))
		}
	}
}
