package daemon

import (
	"strconv"
	"time"

	"github.com/HanksJCTsai/goidleguard/internal/config"
	"github.com/HanksJCTsai/goidleguard/internal/preventidle"
	"github.com/HanksJCTsai/goidleguard/internal/schedule"
	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

type Controller struct {
	cfg        *config.APPConfig
	idleCtl    *preventidle.IdleController
	scheduler  *schedule.Scheduler
	healthStop chan struct{}
}

func NewController(cfg *config.APPConfig) *Controller {
	return &Controller{
		cfg:        cfg,
		idleCtl:    preventidle.NewIdleController(),
		scheduler:  schedule.InitialScheduler(cfg),
		healthStop: make(chan struct{}),
	}
}

func (c *Controller) StartDaemon() {
	logger.LogInfo("Starting daemon...")
	// 啟動持續模擬輸入
	c.idleCtl.StartIdlePrevention()

	task := func() {
		now := time.Now()
		if !c.scheduler.CheckWorkTime(now) {
			if err := preventidle.SimulateActivity(); err != nil {
				logger.LogError("Scheduled SimulateActivity error:", err)
			} else {
				logger.LogInfo("Scheduled activity simulated")
			}
		}
	}

	c.scheduler.ScheduleTask(task)
	// 啟動健康檢查
	// go c.healthStop
}

func (c *Controller) StopDaemon() {
	logger.LogInfo("Stopping daemon...")
	// 停健康檢查
	close(c.healthStop)
	// 停排程與持續輸入模擬
	c.scheduler.StopScheduler()
	c.idleCtl.StopIdlePrevention()
}

func (c *Controller) RestartDaemon() {
	logger.LogInfo("Restarting daemon...")
	c.StopDaemon()
	// 確保資源釋放
	time.Sleep(100 * time.Millisecond)
	c.healthStop = make(chan struct{})
	c.scheduler = schedule.InitialScheduler(c.cfg)
	c.StartDaemon()
}

func (c *Controller) healthCheckLoop() {
	scheduler_interval, err := strconv.Atoi(c.cfg.Scheduler.Interval)
	idle_interval, err := strconv.Atoi(c.cfg.IdlePrevention.Interval)
	if err != nil {
		logger.LogError("invalid idle prevention interval: %w", err)
		return
	}
	if scheduler_interval <= 0 {
		logger.LogError("idle prevention interval must be positive")
	}

	ticker := time.NewTicker(time.Duration(scheduler_interval) * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-c.healthStop:
			logger.LogInfo("Health check stopped")
			return
		case <-ticker.C:
			idleTime, err := preventidle.GetIdleTime()
			if err != nil {
				logger.LogError("HealthCheck: failed to get idle time:", err)
				continue
			}
			// 如果閒置時間過長（例如 10 分鐘以上），可能代表模擬失效，嘗試重啟
			if idleTime > time.Duration(idle_interval+1)*time.Minute {
				logger.LogError("HealthCheck: idle time too long (", idleTime, "), restarting prevention")
				c.RestartDaemon()
			} else {
				logger.LogInfo("HealthCheck: idle time healthy (", idleTime, ")")
			}
		}
	}
}
