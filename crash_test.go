package crash

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitPanicFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "panic.log")

	err := InitPanicFile(tmpFile)
	if err != nil {
		t.Fatalf("InitPanicFile failed: %v", err)
	}

	// 写入一些内容到 stderr
	fmt.Fprintln(os.Stderr, "test panic message")

	// 稍等一下确保写入完成
	time.Sleep(100 * time.Millisecond)

	// 验证文件存在且有内容
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read panic file: %v", err)
	}

	if len(content) == 0 {
		t.Error("panic file is empty")
	}

	t.Logf("Panic file content: %s", content)
}

func TestInitPanicFileWithTee(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "panic_tee.log")

	err := InitPanicFileWithTee(tmpFile)
	if err != nil {
		t.Fatalf("InitPanicFileWithTee failed: %v", err)
	}

	// 写入一些内容到 stderr
	testMsg := "test tee panic message"
	fmt.Fprintln(os.Stderr, testMsg)

	// 稍等一下确保写入完成
	time.Sleep(100 * time.Millisecond)

	// 验证文件存在且有内容
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read panic file: %v", err)
	}

	if len(content) == 0 {
		t.Error("panic file is empty")
	}

	t.Logf("Panic file content: %s", content)
}

func TestPanicHandler(t *testing.T) {
	var handlerCalled bool
	var capturedPanic interface{}

	// 注册一个处理器
	AddPanicHandler(func(panicValue interface{}, stackTrace []byte) {
		handlerCalled = true
		capturedPanic = panicValue
	})

	// 清理
	defer ClearPanicHandlers()

	// 在一个函数中触发 panic 并恢复
	func() {
		defer Recover(false)
		panic("test panic")
	}()

	if !handlerCalled {
		t.Error("panic handler was not called")
	}

	if capturedPanic != "test panic" {
		t.Errorf("expected panic value 'test panic', got: %v", capturedPanic)
	}
}

func TestRecoverToFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "recover.log")

	// 注册文件处理器
	AddPanicHandler(RecoverToFile(tmpFile))
	defer ClearPanicHandlers()

	// 在一个函数中触发 panic 并恢复
	func() {
		defer Recover(false)
		panic("test panic to file")
	}()

	// 稍等一下确保写入完成
	time.Sleep(100 * time.Millisecond)

	// 验证文件存在且有内容
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read recover file: %v", err)
	}

	if len(content) == 0 {
		t.Error("recover file is empty")
	}

	t.Logf("Recover file content: %s", content)
}

func TestWrapMain(t *testing.T) {
	var executed bool

	WrapMain(func() {
		executed = true
	})

	if !executed {
		t.Error("wrapped main function was not executed")
	}
}

func TestRedirectLog(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "app.log")

	err := RedirectLog(tmpFile)
	if err != nil {
		t.Fatalf("RedirectLog failed: %v", err)
	}

	// 写入一些日志（使用 log 包，不是 fmt）
	log.Println("test log message 1")
	log.Println("test log message 2")

	// 稍等一下确保写入完成
	time.Sleep(100 * time.Millisecond)

	// 验证文件存在且有内容
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("log file is empty")
	}

	t.Logf("Log file content: %s", content)
}
