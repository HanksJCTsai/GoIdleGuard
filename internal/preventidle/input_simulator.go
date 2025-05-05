package preventidle

import (
	"fmt"
	"strings"

	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

func SimulateActivity(mode string) error {

	var actions []SimulateAction
	switch strings.ToUpper(mode) {
	default:
	case "MIXED":
		actions = []SimulateAction{
			{"key", "key press"},
			{"mouse", "mouse move"},
		}
	case "KEY":
		actions = []SimulateAction{
			{"key", "key press"},
		}
	case "MOUSE":
		actions = []SimulateAction{
			{"mouse", "mouse move"},
		}
	}

	for _, a := range actions {
		if err := CallSendInput(a.inputType); err != nil {
			return fmt.Errorf("simulate %s failed: %w", a.actionName, err)
		}
		logger.LogInfo("Simulated %s", a.actionName)
	}
	logger.LogInfo("Simulated combined activity")
	return nil
}
