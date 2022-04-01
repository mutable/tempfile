//go:build linux
// +build linux

package tempfile_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mutable/tempfile"
)

func TestOpenCommit(t *testing.T) {
	tmp := os.TempDir()
	pid := os.Getpid()

	f, err := tempfile.Open(tmp, fmt.Sprintf("hello-%d", pid), 0600)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { os.Remove(f.Name()) })

	if err := tempfile.Commit(f); err != nil {
		t.Fatal(err)
	}
}
