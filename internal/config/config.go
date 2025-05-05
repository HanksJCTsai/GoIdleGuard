package config

import (
	"fmt"
	"os"
	"time"
)

// LoadConfig 讀取指定檔案（例如 config.yaml），並反序列化成 Config 結構。
// 同時會呼叫 ValidateConfig 進行設定驗證。
func LoadConfig(path string) (*APPConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseYAMLConfig(data)
	if err != nil {
		return nil, err
	}

	// 驗證設定內容
	if err := ValidateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig 將 cfg 序列化後，原子性地寫入指定檔案。
// 寫入前先寫入暫存檔，再 rename 到正式檔案。
func SaveConfig(path string, cfg *APPConfig) error {
	data, err := MarshalYAML(cfg)
	if err != nil {
		return err
	}

	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, path)
}

// ValidateConfig 驗證設定檔中的各欄位格式與範圍是否正確。
func ValidateConfig(cfg *APPConfig) error {
	if cfg.Scheduler.Interval <= 0 {
		return fmt.Errorf("invalid Scheduler.interval must be >0 (%s)", cfg.Scheduler.Interval)
	}

	// 驗證 IdlePrevention 的 Interval 格式
	if cfg.IdlePrevention.Interval <= 0 {
		return fmt.Errorf("invalid idlePrevention.interval must be >0 (%s)", cfg.IdlePrevention.Interval)
	}

	if cfg.Scheduler.Interval >= cfg.IdlePrevention.Interval {
		return fmt.Errorf("IdlePrevention.tick (%v) must be <= IdlePrevention.interval (%v)",
			cfg.Scheduler.Interval, cfg.IdlePrevention.Interval)
	}

	// 驗證 IdlePrevention 的 Mode 值是否正確
	if cfg.IdlePrevention.Mode != "key" &&
		cfg.IdlePrevention.Mode != "mouse" &&
		cfg.IdlePrevention.Mode != "mixed" {
		return errInvalidMode
	}

	// 驗證 RetryPolicy 的 RetryInterval 格式
	if _, err := time.ParseDuration(cfg.RetryPolicy.RetryInterval); err != nil {
		return fmt.Errorf("invalid retryPolicy.retryInterval format (%s): %w", cfg.RetryPolicy.RetryInterval, err)
	}

	// 驗證 WorkSchedule 每日的工作時段
	for day, sessions := range cfg.WorkSchedule {
		for _, session := range sessions {
			start, err := time.Parse("15:04", session.Start)
			if err != nil {
				return fmt.Errorf("invalid workSchedule.%s start time (%s): %w", day, session.Start, err)
			}
			end, err := time.Parse("15:04", session.End)
			if err != nil {
				return fmt.Errorf("invalid workSchedule.%s end time (%s): %w", day, session.End, err)
			}
			if !start.Before(end) {
				return fmt.Errorf("in workSchedule for %s, start time (%s) must be before end time (%s)", day, session.Start, session.End)
			}
		}
	}

	return nil
}

var errInvalidMode = &InvalidModeError{"Invalid idle prevention mode; must be one of: key, mouse, mixed"}

func (e *InvalidModeError) Error() string {
	return e.Message
}
