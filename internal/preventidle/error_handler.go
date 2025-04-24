package preventidle

import (
	"fmt"

	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

// HandleError 統一處理防止閒置過程中發生的錯誤，這裡範例僅記錄錯誤訊息。
func HandleError(err error) {
	// 可以根據錯誤類型做進一步的重試或通知處理
	msg := fmt.Sprintf("Idle prevention error: %v", err)
	logger.LogError(msg)
}
