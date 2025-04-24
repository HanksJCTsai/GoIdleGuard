package schedule

import (
	"sync"

	"github.com/HanksJCTsai/goidleguard/internal/config"
)

type Scheduler struct {
	Config   *config.APPConfig
	StopChan chan struct{}
	WG       sync.WaitGroup
}
