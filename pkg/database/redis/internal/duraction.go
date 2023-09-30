package internal

import "time"

func Min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
