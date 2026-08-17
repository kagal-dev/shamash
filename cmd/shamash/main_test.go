package main

import (
	"bytes"
	"testing"

	"darvaza.org/core"
)

// TestRootRuns is a boot check: the root command executes and prints.
func TestRootRuns(t *testing.T) {
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetArgs(nil)

	core.AssertNoError(t, rootCmd.Execute(), "execute")
	core.AssertTrue(t, out.Len() > 0, "output")
}
