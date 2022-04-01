// +build linux

package tempfile

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/unix"
)

// Open returns an anonymous file in the directory specified by dir,
// with filepath.Join(dir, name) as file.Name().
func Open(dir, name string, mode os.FileMode) (*os.File, error) {
	if extraBits := mode &^ mode.Perm(); extraBits != 0 {
		return nil, fmt.Errorf("tempfile: unsupported mode bits: %o", extraBits)
	}
	fd, err := unix.Open(dir, unix.O_CLOEXEC|unix.O_RDWR|unix.O_TMPFILE, uint32(mode.Perm()))
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: dir, Err: err}
	}
	return os.NewFile(uintptr(fd), filepath.Join(dir, name)), nil
}

// Commit links a file returned by Open into its directory
// after invoking f.Sync().
func Commit(f *os.File) error {
	if err := f.Sync(); err != nil {
		return err
	}

	err := unix.Linkat(
		unix.AT_FDCWD, fmt.Sprintf("/proc/self/fd/%d", f.Fd()),
		unix.AT_FDCWD, f.Name(),
		unix.AT_SYMLINK_FOLLOW,
	)
	runtime.KeepAlive(f)

	if err != nil {
		return &os.PathError{Op: "linkat", Path: f.Name(), Err: err}
	}
	return nil
}
