package daemon

import (
	"testing"
	"time"

	"github.com/HanksJCTsai/goidleguard/internal/config"
	"github.com/HanksJCTsai/goidleguard/internal/preventidle"
	"github.com/HanksJCTsai/goidleguard/internal/schedule"
)

// 整合測試：使用真實 Controller + 真實模組，來測試是否能成功啟動與停止
func TestDaemonController_StartAndStop(t *testing.T) {
	cfg := &config.APPConfig{
		Version: config.VersionConfig{Name: "TestApp",
			Version: "0.1.0"},
		Scheduler: config.SchedulerConfig{Interval: "5m"},
		IdlePrevention: config.IdlePreventionConfig{
			Enabled:  true,
			Interval: "1s", // 縮短執行間隔，方便測試
			Mode:     "key",
		},
		Logging: config.LoggingConfig{
			Level:  "debug",
			Output: "console",
		},
		RetryPolicy: config.RetryPolicyConfig{
			RetryInterval: "1s",
			MaxRetries:    1,
		},
		WorkSchedule: config.WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
			"tuesday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
			"wednesday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
			"thursday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
			"friday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
			"saturday": {},
			"sunday":   {},
		},
	}

	// 建立真實的 Controller 實例
	ctrl := &Controller{
		cfg:        cfg,
		idleCtl:    preventidle.NewIdleController(),
		scheduler:  schedule.InitialScheduler(cfg),
		healthStop: make(chan struct{}),
	}

	// 啟動 Daemon
	ctrl.StartDaemon()
	t.Log("→ Daemon started")

	// 等待 3 秒鐘，讓背景排程模組有機會跑起來
	time.Sleep(3 * time.Second)

	// 停止 Daemon
	ctrl.StopDaemon()
	t.Log("   → Daemon stopped")

	// 如果沒有 panic、沒有錯誤，代表啟動與停止都正常
}
