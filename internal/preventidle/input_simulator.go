package preventidle

import (
	"fmt"

	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

func SimulateKeyPress() error {
	if err := CallSendInput("key"); err != nil {
		return fmt.Errorf("SimulateKeyPress failed: %w", err)
	}
	logger.LogInfo("Simulated key press")
	return nil
}

func SimulateMouseMove() error {
	if err := CallSendInput("mouse"); err != nil {
		return fmt.Errorf("SimulateMouseMove failed: %w", err)
	}
	logger.LogInfo("Simulated mouse move")
	return nil
}

func SimulateActivity() error {
	// if err := SimulateKeyPress(); err != nil {
	// 	return err
	// }
	if err := SimulateMouseMove(); err != nil {
		return err
	}
	logger.LogInfo("Simulated combined activity")
	return nil
}
