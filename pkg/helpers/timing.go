package helpers

import "time"

func Timing(cb func()) time.Duration {
	ts := time.Now()
	cb()
	return time.Since(ts)
}
