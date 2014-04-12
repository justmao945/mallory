package mallory

import (
	"net/http"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

func CopyResponseHeader(w http.ResponseWriter, r *http.Response) {
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

// From https://github.com/mreiferson/go-simplelog/blob/master/simplelog.go
func Ioctl(fd, request, argp uintptr) syscall.Errno {
	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, fd, request, argp)
	return errorp
}

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
