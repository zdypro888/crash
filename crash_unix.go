// +build darwin linux

package crash

import (
	"os"
	"syscall"
)

//InitPanicFile 初始化 crash 文件
func InitPanicFile(panicFile string) (*os.File, error) {
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	return file, nil
}
