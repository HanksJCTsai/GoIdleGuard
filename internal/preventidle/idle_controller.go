package preventidle

import (
	"time"

	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

func NewIdleController() *IdleController {
	return &IdleController{
		StopChan: make(chan struct{}),
		Running:  false,
	}
}

func (ic *IdleController) StartIdlePrevention() {
	if ic.Running {
		logger.LogInfo("Idle prevention already running")
		return
	}
	ic.StopChan = make(chan struct{})
	ic.Running = true
	go ic.run()
	logger.LogInfo("Idle prevention started")
}

func (ic *IdleController) run() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ic.StopChan:
			logger.LogInfo("Idle prevention stopped")
			return
		case <-ticker.C:
			if err := SimulateActivity(); err != nil {
				HandleError(err)
			}
		}
	}
}

func (ic *IdleController) StopIdlePrevention() {
	if ic.Running == false {
		logger.LogInfo("Idle prevention is not running")
		return
	}
	close(ic.StopChan)
	ic.Running = false
	logger.LogInfo("Idle prevention stop requested")
}

func (ic *IdleController) RestartIdlePrevention() {
	ic.StopIdlePrevention()
	time.Sleep(1000 * time.Microsecond)
	ic.StartIdlePrevention()
}
