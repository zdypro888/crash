// +build darwin linux

package crash

import (
	"os"
	"runtime"
	"syscall"
)

//InitPanicFile 初始化 crash 文件
func InitPanicFile(panicFile string) error {
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		file.Close()
		return err
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(file, func(fd *os.File) {
		fd.Close()
	})
	return nil
}
