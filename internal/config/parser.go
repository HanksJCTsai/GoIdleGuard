package config

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// ParseYAMLConfig 解析 YAML 格式的資料成為 Config 結構
func ParseYAMLConfig(data []byte) (*APPConfig, error) {
	var cfg APPConfig
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MarshalYAML 將 Config 結構序列化成 YAML 格式資料
func MarshalYAML(cfg *APPConfig) ([]byte, error) {
	return yaml.Marshal(cfg)
}

// ParseJSONConfig 解析 JSON 格式的資料成為 Config 結構
func ParseJSONConfig(data []byte) (*APPConfig, error) {
	var cfg APPConfig
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MarshalJSON 將 Config 結構序列化成 JSON 格式資料
func MarshalJSON(cfg *APPConfig) ([]byte, error) {
	return json.MarshalIndent(cfg, "", "  ")
}
