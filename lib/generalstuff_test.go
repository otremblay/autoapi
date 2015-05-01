package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoTablesDoNotPanic(t *testing.T) {
	dontPanicBro(t, &checksumGenerator{})
	dontPanicBro(t, &dbCodeGenerator{})
	dontPanicBro(t, &dbiCodeGenerator{})
	dontPanicBro(t, &httpCodeGenerator{})
	dontPanicBro(t, &swaggerGenerator{})
	dontPanicBro(t, &halgenerator{})
	dontPanicBro(t, &mainGenerator{})
}
func dontPanicBro(t *testing.T, g Generator) {
	fn := func() {
		g.Generate(map[string]tableInfo{})
		g.Generate(nil)
	}
	assert.NotPanics(t, fn)
}
