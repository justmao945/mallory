package mallory

import (
	"net/http"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

// Duration to e.g. 432ms or 12s, human readable translation
func BeautifyDuration(d time.Duration) string {
	u, ms, s := uint64(d), uint64(time.Millisecond), uint64(time.Second)
	if d < 0 {
		u = -u
	}
	switch {
	case u < ms:
		return "0"
	case u < s:
		return strconv.FormatUint(u/ms, 10) + "ms"
	default:
		return strconv.FormatUint(u/s, 10) + "s"
	}
}

// copy and overwrite headers from r to w
func CopyHeader(w http.ResponseWriter, r *http.Response) {
	// copy headers
	dst, src := w.Header(), r.Header
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

// POSIX ioctl syscall
// From https://github.com/mreiferson/go-simplelog/blob/master/simplelog.go
func Ioctl(fd, request, argp uintptr) syscall.Errno {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, fd, request, argp)
	return errorp
}

// Test is a termnial or not
// From https://github.com/mreiferson/go-simplelog/blob/master/simplelog.go
func Isatty(f *os.File) bool {
	switch runtime.GOOS {
	case "darwin":
	case "linux":
	default:
		return false
	}
	var t [2]byte
	errno := Ioctl(f.Fd(), syscall.TIOCGPGRP, uintptr(unsafe.Pointer(&t)))
	return errno == 0
}
