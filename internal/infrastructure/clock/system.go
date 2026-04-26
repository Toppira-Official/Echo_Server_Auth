package clock

import "time"

type SystemClock struct{}

func NewSystemClock() *SystemClock { return &SystemClock{} }

func (*SystemClock) NowUTC() time.Time {
	return time.Now().UTC()
}
