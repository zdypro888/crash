package crash

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

// PanicHandler 是 panic 发生时的回调函数类型
// panicValue 是 panic 的值，stackTrace 是堆栈信息
type PanicHandler func(panicValue interface{}, stackTrace []byte)

var (
	panicHandlers []PanicHandler
	handlerMutex  sync.RWMutex
)

// AddPanicHandler 添加一个 panic 处理器
// 当 panic 发生时，所有已注册的处理器都会被调用
func AddPanicHandler(handler PanicHandler) {
	handlerMutex.Lock()
	defer handlerMutex.Unlock()
	panicHandlers = append(panicHandlers, handler)
}

// ClearPanicHandlers 清除所有已注册的 panic 处理器
func ClearPanicHandlers() {
	handlerMutex.Lock()
	defer handlerMutex.Unlock()
	panicHandlers = nil
}

// Recover 应该在 defer 中调用，用于捕获 panic 并执行已注册的处理器
// 如果 rePanic 为 true，则在处理完成后会重新 panic
func Recover(rePanic bool) {
	if r := recover(); r != nil {
		stack := debug.Stack()

		// 调用所有已注册的处理器
		handlerMutex.RLock()
		handlers := make([]PanicHandler, len(panicHandlers))
		copy(handlers, panicHandlers)
		handlerMutex.RUnlock()

		for _, handler := range handlers {
			func() {
				defer func() {
					// 防止处理器本身 panic
					if err := recover(); err != nil {
						fmt.Fprintf(os.Stderr, "panic handler failed: %v\n", err)
					}
				}()
				handler(r, stack)
			}()
		}

		if rePanic {
			panic(r)
		}
	}
}

// WrapMain 包装 main 函数，自动捕获并处理 panic
// 示例：
//
//	func main() {
//	    crash.WrapMain(func() {
//	        // 你的 main 函数逻辑
//	    })
//	}
func WrapMain(mainFunc func()) {
	defer Recover(false)
	mainFunc()
}

// RecoverToFile 是一个便捷函数，将 panic 信息写入指定文件
// 这个函数可以直接作为 PanicHandler 使用
func RecoverToFile(filename string) PanicHandler {
	return func(panicValue interface{}, stackTrace []byte) {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open panic file: %v\n", err)
			return
		}
		defer f.Close()

		fmt.Fprintf(f, "\n=== PANIC ===\n")
		fmt.Fprintf(f, "Value: %v\n", panicValue)
		fmt.Fprintf(f, "Stack Trace:\n%s\n", stackTrace)
	}
}
