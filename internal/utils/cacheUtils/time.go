package cacheUtils

import (
	"time"
)

func Seconds(duration int, durationType time.Duration) int64 {
	return int64((time.Duration(duration) * durationType).Seconds())
}
