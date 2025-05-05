package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/HanksJCTsai/goidleguard/internal/config"
	"github.com/HanksJCTsai/goidleguard/pkg/logger"
)

func main() {
	logger.InitLogger()
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.LogError("Failed to load config:", err)
		os.Exit(1)
	}
	logger.LogInfo("Config loaded successfully")

	// 建立並啟動 DaemonController
	dc := NewController(cfg)
	dc.StartDaemon()

	// 捕捉系統中斷訊號以優雅關閉
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.LogInfo("Shutdown signal received, stopping daemon...")
	dc.StopDaemon()
	logger.LogInfo("Daemon stopped; exiting.")
}
