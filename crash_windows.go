// +build windows

package crash

import (
	"os"
	"syscall"
)

const (
	kernel32dll = "kernel32.dll"
)

//InitPanicFile 初始化 crash 文件
func InitPanicFile(panicFile string) (*os.File, error) {
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	kernel32 := syscall.NewLazyDLL(kernel32dll)
	setStdHandle := kernel32.NewProc("SetStdHandle")
	sh := syscall.STD_ERROR_HANDLE
	v, _, err := setStdHandle.Call(uintptr(sh), uintptr(file.Fd()))
	if v == 0 {
		return nil, err
	}
	return file, nil
}
