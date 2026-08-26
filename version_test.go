package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintVersion_BuildInfoNotAvailable(t *testing.T) {
	t.Parallel()

	// When not built with VCS info (e.g., go run without git),
	// it should print "version unknown"
	var buf bytes.Buffer
	printVersion(&buf)

	out := buf.String()
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "dl")
}
