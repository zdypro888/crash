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
func InitPanicFile(panicFile string) error {
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	kernel32 := syscall.NewLazyDLL(kernel32dll)
	setStdHandle := kernel32.NewProc("SetStdHandle")
	sh := syscall.STD_ERROR_HANDLE
	v, _, err := setStdHandle.Call(uintptr(sh), uintptr(file.Fd()))
	if v == 0 {
		fd.Close()
		return err
	}
	// runtime.SetFinalizer(file, func(fd *os.File) {
	// 	fd.Close()
	// })
	return nil
}
