package schedule

import (
	"time"

	"github.com/HanksJCTsai/goidleguard/internal/config"
)

func parseSessionTime(tStr string, now time.Time) (time.Time, error) {
	layout := "15:04"
	paresd, err := time.ParseInLocation(layout, tStr, now.Location())
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(now.Year(), now.Month(), now.Day(), paresd.Hour(), paresd.Minute(), 0, 0, now.Location()), nil
}

// IsTimeInRange 判斷 target 是否介於 start 與 end 之間。
func IsTimeInRange(target, start, end time.Time) bool {
	return target.After(start) && target.Before(end)
}

func CalculateNextInterval(cfg *config.APPConfig) (time.Duration, error) {
	duration, err := time.ParseDuration(cfg.Scheduler.Interval)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
