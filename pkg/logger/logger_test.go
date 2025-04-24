package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestInitLogger(t *testing.T) {
	// 呼叫 InitLogger，檢查內部 loggers 是否已初始化
	InitLogger()
	if infoLogger == nil || debugLogger == nil || errorLogger == nil {
		t.Error("Expected loggers to be initialized")
	}
}

func TestLogFunctions(t *testing.T) {
	// 將 logger 輸出導向 buffer，以便檢查輸出內容
	var buf bytes.Buffer
	infoLogger = log.New(&buf, "INFO: ", 0)
	debugLogger = log.New(&buf, "DEBUG: ", 0)
	errorLogger = log.New(&buf, "ERROR: ", 0)

	LogInfo("Test Info")
	LogDebug("Test Debug")
	LogError("Test Error")

	output := buf.String()
	if !contains(output, "Test Info") || !contains(output, "Test Debug") || !contains(output, "Test Error") {
		t.Errorf("Log functions did not produce expected output: %s", output)
	}
}

// contains 是檢查子字串是否存在的輔助函式
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
