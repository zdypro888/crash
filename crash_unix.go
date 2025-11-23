//go:build darwin || linux

package crash

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

//InitPanicFile 初始化 crash 文件，将 stderr 重定向到文件
func InitPanicFile(panicFile string) error {
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if err = unix.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		file.Close()
		return err
	}
	// 内存回收前关闭文件描述符
	// runtime.SetFinalizer(file, func(fd *os.File) {
	// 	fd.Close()
	// })
	return nil
}

//InitPanicFileWithTee 初始化 crash 文件，同时输出到控制台和文件
func InitPanicFileWithTee(panicFile string) error {
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 保存原始的 stderr
	origStderr, err := unix.Dup(int(os.Stderr.Fd()))
	if err != nil {
		file.Close()
		return err
	}

	// 创建管道
	r, w, err := os.Pipe()
	if err != nil {
		unix.Close(origStderr)
		file.Close()
		return err
	}

	// 将 stderr 重定向到管道的写端
	if err = unix.Dup2(int(w.Fd()), int(os.Stderr.Fd())); err != nil {
		r.Close()
		w.Close()
		unix.Close(origStderr)
		file.Close()
		return err
	}

	// 启动 goroutine 将管道读端的内容同时写入文件和原始 stderr
	go func() {
		origFile := os.NewFile(uintptr(origStderr), "/dev/stderr")
		mw := io.MultiWriter(origFile, file)
		io.Copy(mw, r)
	}()

	return nil
}
