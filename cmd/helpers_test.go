package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIfError(t *testing.T) {
	// Test for checkIfError function.

	t.Run("No error", func(t *testing.T) {
		// Test when err is nil.  We expect no output and no exit.
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Call checkIfError with a nil error.
		checkIfError(nil)

		// Check that nothing was written to standard output.
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()
		assert.Empty(t, output, "Should not print anything when err is nil")
	})

}
