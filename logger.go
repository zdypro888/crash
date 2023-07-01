package crash

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

func RedirectLog(file string) error {
	// logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0644)
	// if err != nil {
	// 	return err
	// }
	// runtime.SetFinalizer(logFile, func(fd *os.File) {
	// 	fd.Close()
	// })
	// 创建一个 lumberjack.Logger 实例，用于管理日志文件
	logFile := &lumberjack.Logger{
		Filename:   file, // 日志文件名
		MaxSize:    100,  // 日志文件的最大大小（以MB为单位）
		MaxBackups: 3,    // 最多保留的旧日志文件数
		MaxAge:     30,   // 旧日志文件保留的最长天数
		Compress:   true, // 是否启用压缩旧日志文件
	}
	log.SetOutput(logFile)
	return nil
}
