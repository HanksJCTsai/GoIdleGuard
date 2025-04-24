package config

// Config 定義了從 config.yaml 讀取的整個設定結構
type APPConfig struct {
	Version        VersionConfig        `yaml:"version" json:"version"`
	Scheduler      SchedulerConfig      `yaml:"scheduler" json:"scheduler"`
	IdlePrevention IdlePreventionConfig `yaml:"idlePrevention" json:"idlePrevention"`
	Logging        LoggingConfig        `yaml:"logging" json:"logging"`
	RetryPolicy    RetryPolicyConfig    `yaml:"retryPolicy" json:"retryPolicy"`
	WorkSchedule   WorkSchedule         `yaml:"workSchedule" json:"workSchedule"`
}

type VersionConfig struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

type IdlePreventionConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Interval string `yaml:"interval" json:"interval"` // 例如 "5m"
	Mode     string `yaml:"mode" json:"mode"`         // 可選值： "key"、"mouse"、"mixed"
}

type SchedulerConfig struct {
	Interval string `yaml:"interval" json:"interval"` // 例如 "10m"
}

type LoggingConfig struct {
	Level  string `yaml:"level" json:"level"`   // 例如 "info" 或 "debug"
	Output string `yaml:"output" json:"output"` // "console" 或檔案路徑
}

type RetryPolicyConfig struct {
	MaxRetries    int    `yaml:"maxRetries" json:"maxRetries"`
	RetryInterval string `yaml:"retryInterval" json:"retryInterval"` // 例如 "10s"
}

type InvalidModeError struct {
	Message string
}

// WorkSession 定義一天內單個工作時段的開始與結束時間
type WorkSession struct {
	Start string `yaml:"start" json:"start"`
	End   string `yaml:"end" json:"end"`
}

// WorkSchedule 定義一週內每天的工作時段，使用 map 對應每一天的時段陣列
type WorkSchedule map[string][]WorkSession
