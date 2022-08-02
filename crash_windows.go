//go:build windows

package crash

import (
	"golang.org/x/sys/windows"
	"os"
)

//InitPanicFile 初始化 crash 文件
func InitPanicFile(panicFile string) error {
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	return windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(file.Fd()))
}
