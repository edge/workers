package workers

import "time"

func GetScheduledTicker(t time.Time) <-chan time.Time {
	now := time.Now()
	diff := t.Sub(now)
	return time.After(diff)
}
