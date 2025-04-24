package schedule

import (
	"strings"
	"time"

	"github.com/HanksJCTsai/goidleguard/internal/config"
)

func InitialScheduler(cfg *config.APPConfig) *Scheduler {
	return &Scheduler{
		Config:   cfg,
		StopChan: make(chan struct{}),
	}
}

func (s *Scheduler) CheckWorkTime(now time.Time) bool {
	day := strings.ToLower(now.Weekday().String()) // 例如 "monday"
	sessions, exists := s.Config.WorkSchedule[day]
	if !exists || len(sessions) == 0 {
		return false
	}
	for _, session := range sessions {
		start, err := parseSessionTime(session.Start, now)
		if err != nil {
			continue
		}
		end, err := parseSessionTime(session.End, now)
		if err != nil {
			continue
		}
		if IsTimeInRange(now, start, end) {
			return true
		}
	}
	return false
}

func (s *Scheduler) ScheduleTask(task func()) {
	s.WG.Add(1)
	go func() {
		defer s.WG.Done()
		var now time.Time
		for {
			select {
			case <-s.StopChan:
				return
			default:
				now = time.Now()
				if !s.CheckWorkTime(now) {
					task()
				}

				interval, err := CalculateNextInterval(s.Config)
				if err != nil {
					return
				}

				timer := time.NewTimer(interval)
				select {
				case <-timer.C:
					// 繼續下一輪任務檢查
				case <-s.StopChan:
					timer.Stop()
					return
				}
			}
		}
	}()
}

func (s *Scheduler) StopScheduler() {
	close(s.StopChan)
	s.WG.Wait()
}
