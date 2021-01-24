package types

import (
	"time"
)

type VObj interface {
	GetAsFloat() (float64, error)
	GetAsInt() (int64, error)
	GetAsDatetime() (time.Time, error)
	GetAsPercent() (float64, error)
}
