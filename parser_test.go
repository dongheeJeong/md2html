package main

import "testing"

func TestReplacePrefix(t *testing.T) {

  got := ReplacePrefix("# abc", "# ", "<h1>")
  want := "<h1> abc"

  if got != want {
    t.Errorf("want %q, got %q", want, got)
  }
}
