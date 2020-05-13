package render

import (
	"testing"

	"github.com/sh0rez/docsonnet/pkg/docsonnet"
	"github.com/stretchr/testify/assert"
)

func TestSortFields(t *testing.T) {
	api := docsonnet.Fields{
		"new":      dfn(),
		"newNamed": dfn(),

		"aaa": dfn(),
		"bbb": dobj(),
		"ccc": dfn(),

		"metadata": dobj(),
	}

	sorted := []string{
		"new",
		"newNamed",

		"aaa",
		"ccc",

		"bbb",
		"metadata",
	}

	res := sortFields(api)

	assert.Equal(t, sorted, res)
}

func dobj() docsonnet.Field {
	return docsonnet.Field{
		Object: &docsonnet.Object{},
	}
}

func dfn() docsonnet.Field {
	return docsonnet.Field{
		Function: &docsonnet.Function{},
	}
}
