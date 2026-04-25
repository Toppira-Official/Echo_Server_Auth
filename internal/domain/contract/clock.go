package contract

import "time"

type Clock interface {
	Now_UTC() time.Time
}
