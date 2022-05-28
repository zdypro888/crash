package crash

import (
	"log"
	"os"
	"runtime"
)

func RedirectLog(file string) error {
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	runtime.SetFinalizer(logFile, func(fd *os.File) {
		fd.Close()
	})
	log.SetOutput(logFile)
	return nil
}
