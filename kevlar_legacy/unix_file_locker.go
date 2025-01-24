package kevlar_legacy

import (
	"io"
	"syscall"
)

// TODO: remove this when https://github.com/golang/go/issues/33974 is accepted

func lockFd(fd uintptr) error {
	return syscall.FcntlFlock(fd, syscall.F_SETLK, &syscall.Flock_t{
		Start:  0,
		Len:    0,
		Type:   syscall.F_RDLCK | syscall.F_WRLCK,
		Whence: io.SeekStart,
	})
}

func unlockFd(fd uintptr) error {
	return syscall.FcntlFlock(fd, syscall.F_SETLK, &syscall.Flock_t{
		Start:  0,
		Len:    0,
		Type:   syscall.F_UNLCK,
		Whence: io.SeekStart,
	})
}
