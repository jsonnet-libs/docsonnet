package md

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	l := List(
		Text("foo"),
		Text("bar"),
		List(
			Text("baz"),
			Text("bing"),
		),
		Text("boing"),
	).String()

	assert.Equal(t, `* foo
* bar
  * baz
  * bing
* boing`, l)
}
