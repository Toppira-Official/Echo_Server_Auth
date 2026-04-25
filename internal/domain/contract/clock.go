package contract

import "time"

type Clock interface {
	NowUTC() time.Time
}
