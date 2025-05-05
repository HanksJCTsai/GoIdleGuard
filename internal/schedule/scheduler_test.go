package schedule

import (
	"testing"
	"time"

	"github.com/HanksJCTsai/goidleguard/internal/config"
)

func TestParseTimeString(t *testing.T) {
	ref := time.Date(2025, time.April, 2, 0, 0, 0, 0, time.Local)
	tStr := "09:30"
	parsed, err := parseSessionTime(tStr, ref)
	if err != nil {
		t.Fatalf("ParseTimeString failed: %v", err)
	}
	expected := time.Date(2025, time.April, 2, 9, 30, 0, 0, time.Local)
	if !parsed.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}

func TestIsTimeInRange(t *testing.T) {
	now := time.Now()
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	if !IsTimeInRange(now, start, end) {
		t.Errorf("Expected now to be in range")
	}
}

func TestCheckWorkTime(t *testing.T) {
	// 建立一組包含 monday 工作時段的配置
	cfg := &config.APPConfig{
		WorkSchedule: config.WorkSchedule{
			"monday": {
				{Start: "08:00", End: "12:00"},
				{Start: "13:00", End: "17:00"},
			},
		},
	}

	// 建立 Scheduler
	s := InitialScheduler(cfg)

	timeInWork := time.Date(2025, time.April, 7, 9, 0, 0, 0, time.Local)
	if !s.CheckWorkTime(timeInWork) {
		t.Errorf("Expected %v to be within work session", timeInWork)
	}

	timeOutWork := time.Date(2025, time.April, 7, 12, 30, 0, 0, time.Local)
	if s.CheckWorkTime(timeOutWork) {
		t.Errorf("Expected %v to be outside work session", timeOutWork)
	}
}

func TestSchedulerScheduleTask(t *testing.T) {
	// 測試 ScheduleTask 是否正確啟動 task
	cfg := &config.APPConfig{
		WorkSchedule: config.WorkSchedule{
			"monday": {
				{Start: "00:00", End: "00:01"}, // 測試用，設定為幾乎不影響
			},
		},
	}
	s := InitialScheduler(cfg)
	done := make(chan bool, 1)
	task := func() {
		done <- true
	}
	// 執行 ScheduleTask
	s.ScheduleTask(task)

	// 等待 task 執行（固定間隔為 1 分鐘，為測試縮短等待時間，這裡僅檢查是否在合理時間內回傳結果）
	select {
	case <-done:
		// 任務執行成功
	case <-time.After(2 * time.Minute):
		t.Error("Task was not executed within expected time")
	}
	s.StopScheduler()
}
