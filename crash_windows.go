//go:build windows

package crash

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
)

//InitPanicFile 初始化 crash 文件，将 stderr 重定向到文件
func InitPanicFile(panicFile string) error {
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(file.Fd()))
}

//InitPanicFileWithTee 初始化 crash 文件，同时输出到控制台和文件
func InitPanicFileWithTee(panicFile string) error {
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 保存原始的 stderr handle
	origStderr, err := windows.GetStdHandle(windows.STD_ERROR_HANDLE)
	if err != nil {
		file.Close()
		return err
	}

	// 创建管道
	r, w, err := os.Pipe()
	if err != nil {
		file.Close()
		return err
	}

	// 将 stderr 重定向到管道的写端
	if err = windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(w.Fd())); err != nil {
		r.Close()
		w.Close()
		file.Close()
		return err
	}

	// 启动 goroutine 将管道读端的内容同时写入文件和原始 stderr
	go func() {
		origFile := os.NewFile(uintptr(origStderr), "CONOUT$")
		mw := io.MultiWriter(origFile, file)
		io.Copy(mw, r)
	}()

	return nil
}
