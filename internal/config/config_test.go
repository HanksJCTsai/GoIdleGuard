package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// 建立包含 workSchedule 的暫存 YAML 配置內容
	content := []byte(`
version:
  name: PreventIdleApp
  version: "1.0.0"

scheduler:
  interval: "10m"

idlePrevention:
  enabled: true
  interval: "5m"
  mode: "mixed"

logging:
  level: "info"
  output: "console"

retryPolicy:
  maxRetries: 3
  retryInterval: "10s"

workSchedule:
  monday:
    - start: "08:00"
      end: "12:00"
    - start: "13:00"
      end: "17:00"
  tuesday:
    - start: "08:00"
      end: "12:00"
    - start: "13:00"
      end: "17:00"
  wednesday:
    - start: "08:00"
      end: "12:00"
    - start: "13:00"
      end: "17:00"
  thursday:
    - start: "08:00"
      end: "12:00"
    - start: "13:00"
      end: "17:00"
  friday:
    - start: "08:00"
      end: "12:00"
    - start: "13:00"
      end: "17:00"
  saturday: []
  sunday: []
`)
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 檢查 App 欄位
	if cfg.Version.Name != "PreventIdleApp" {
		t.Errorf("Expected App.Name 'PreventIdleApp', got %s", cfg.Version.Name)
	}

	// 檢查 WorkSchedule 的解析結果
	sessions, ok := cfg.WorkSchedule["monday"]
	if !ok {
		t.Errorf("Expected monday workSchedule, but not found")
	}
	if len(sessions) != 2 {
		t.Errorf("Expected 2 sessions for monday, got %d", len(sessions))
	}
	if sessions[0].Start != "08:00" || sessions[0].End != "12:00" {
		t.Errorf("Expected first monday session to be 08:00-12:00, got %s-%s", sessions[0].Start, sessions[0].End)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config_save_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	os.Remove(tmpFile.Name()) // 讓 SaveConfig 自行建立檔案

	cfg := &APPConfig{
		Version: VersionConfig{
			Name:    "TestApp",
			Version: "0.1.0",
		},
		Scheduler: SchedulerConfig{
			Interval: (1 * time.Minute),
		},
		IdlePrevention: IdlePreventionConfig{
			Enabled:  true,
			Interval: (5 * time.Minute),
			Mode:     "key",
		},
		Logging: LoggingConfig{
			Level:  "debug",
			Output: "console",
		},
		RetryPolicy: RetryPolicyConfig{
			MaxRetries:    5,
			RetryInterval: "5s",
		},
		WorkSchedule: WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
			"tuesday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
		},
	}

	if err := SaveConfig(tmpFile.Name(), cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loadedCfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig after SaveConfig failed: %v", err)
	}

	if loadedCfg.Version.Name != cfg.Version.Name {
		t.Errorf("Expected App.Name %s, got %s", cfg.Version.Name, loadedCfg.Version.Name)
	}

	sessions, ok := loadedCfg.WorkSchedule["monday"]
	if !ok || len(sessions) != 2 {
		t.Errorf("Expected 2 sessions for monday, got %v", sessions)
	}
}

func TestValidateConfig_InvalidWorkSchedule(t *testing.T) {
	// 測試當某天工作時段中，開始時間不小於結束時間的情形
	cfg := &APPConfig{
		Version: VersionConfig{
			Name:    "TestApp",
			Version: "0.1.0",
		},
		Scheduler: SchedulerConfig{
			Interval: (1 * time.Minute),
		},
		IdlePrevention: IdlePreventionConfig{
			Enabled:  true,
			Interval: (5 * time.Minute),
			Mode:     "key",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Output: "console",
		},
		RetryPolicy: RetryPolicyConfig{
			MaxRetries:    3,
			RetryInterval: "10s",
		},
		WorkSchedule: WorkSchedule{
			"monday": {
				{Start: "08:00", End: "07:00"}, // 開始時間晚於結束時間
			},
		},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Errorf("Expected error for invalid monday workSchedule, got nil")
	}
}

func TestValidateConfig_InvalidIdlePreventionMode(t *testing.T) {
	cfg := &APPConfig{
		Version: VersionConfig{
			Name:    "TestApp",
			Version: "0.1.0",
		},
		Scheduler: SchedulerConfig{
			Interval: (1 * time.Minute),
		},
		IdlePrevention: IdlePreventionConfig{
			Enabled:  true,
			Interval: (5 * time.Minute),
			Mode:     "invalid_mode",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Output: "console",
		},
		RetryPolicy: RetryPolicyConfig{
			MaxRetries:    3,
			RetryInterval: "10s",
		},
		WorkSchedule: WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
			},
		},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Errorf("Expected error for invalid IdlePrevention.Mode, got nil")
	} else {
		expected := "Invalid idle prevention mode; must be one of: key, mouse, mixed"
		if err.Error() != expected {
			t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
		}
	}
}

func TestValidateConfig_InvalidRetryInterval(t *testing.T) {
	cfg := &APPConfig{
		Version: VersionConfig{
			Name:    "TestApp",
			Version: "0.1.0",
		},
		Scheduler: SchedulerConfig{
			Interval: (1 * time.Minute),
		},
		IdlePrevention: IdlePreventionConfig{
			Enabled:  true,
			Interval: (5 * time.Minute),
			Mode:     "key",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Output: "console",
		},
		RetryPolicy: RetryPolicyConfig{
			MaxRetries:    3,
			RetryInterval: "invalid", // 格式錯誤
		},
		WorkSchedule: WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
			},
		},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Errorf("Expected error for invalid retryPolicy.retryInterval, got nil")
	}
}

func TestValidateConfig_InvalidSchedulerInterval(t *testing.T) {
	cfg := &APPConfig{
		Version: VersionConfig{
			Name:    "TestApp",
			Version: "0.1.0",
		},
		Scheduler: SchedulerConfig{
			Interval: (5 * time.Minute),
		},
		IdlePrevention: IdlePreventionConfig{
			Enabled:  true,
			Interval: (1 * time.Minute), // 格式錯誤
			Mode:     "key",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Output: "console",
		},
		RetryPolicy: RetryPolicyConfig{
			MaxRetries:    3,
			RetryInterval: "10s",
		},
		WorkSchedule: WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
			},
		},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Errorf("Expected error for invalid scheduler.interval, got nil")
	}
}

func TestValidateConfig_InvalidIdlePreventionInterval(t *testing.T) {
	cfg := &APPConfig{
		Version: VersionConfig{
			Name:    "TestApp",
			Version: "0.1.0",
		},
		Scheduler: SchedulerConfig{
			Interval: (5 * time.Minute),
		},
		IdlePrevention: IdlePreventionConfig{
			Enabled:  true,
			Interval: (1 * time.Minute), // 格式錯誤
			Mode:     "key",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Output: "console",
		},
		RetryPolicy: RetryPolicyConfig{
			MaxRetries:    3,
			RetryInterval: "10s",
		},
		WorkSchedule: WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
			},
		},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Errorf("Expected error for invalid idlePrevention.interval, got nil")
	}
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	cfg := &APPConfig{
		Version: VersionConfig{
			Name:    "ValidApp",
			Version: "1.0.0",
		},
		Scheduler: SchedulerConfig{
			Interval: (1 * time.Minute),
		},
		IdlePrevention: IdlePreventionConfig{
			Enabled:  true,
			Interval: (5 * time.Minute),
			Mode:     "mouse",
		},
		Logging: LoggingConfig{
			Level:  "debug",
			Output: "console",
		},
		RetryPolicy: RetryPolicyConfig{
			MaxRetries:    2,
			RetryInterval: "10s",
		},
		WorkSchedule: WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
			"tuesday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
		},
	}

	if err := ValidateConfig(cfg); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestTimeParsingForWorkSchedule(t *testing.T) {
	// 測試工作時段時間格式是否符合 "15:04"
	session := WorkSession{Start: "09:00", End: "17:00"}
	if _, err := time.Parse("15:04", session.Start); err != nil {
		t.Errorf("Failed to parse start time: %v", err)
	}
	if _, err := time.Parse("15:04", session.End); err != nil {
		t.Errorf("Failed to parse end time: %v", err)
	}
}
