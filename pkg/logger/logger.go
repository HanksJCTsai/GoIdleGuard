package logger

import (
	"log"
	"os"
	"sync"
)

var (
	loggerOnce  sync.Once
	infoLogger  *log.Logger
	debugLogger *log.Logger
	errorLogger *log.Logger
)

// InitLogger 初始化 logger，可以配置輸出到標準輸出或檔案。
func InitLogger() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogInfo 輸出 Info 級別的日誌訊息。
func LogInfo(v ...interface{}) {
	loggerOnce.Do(func() {
		InitLogger()
	})
	infoLogger.Println(v...)
}

// LogDebug 輸出 Debug 級別的日誌訊息。
func LogDebug(v ...interface{}) {
	loggerOnce.Do(func() {
		InitLogger()
	})
	debugLogger.Println(v...)
}

// LogError 輸出 Error 級別的日誌訊息。
func LogError(v ...interface{}) {
	loggerOnce.Do(func() {
		InitLogger()
	})
	errorLogger.Println(v...)
}
