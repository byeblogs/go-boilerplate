package constants

import "time"

var (
	NSQRetryDelayMap = map[uint16]time.Duration{
		1: 5 * time.Second,
		2: 1 * time.Minute,
		3: 5 * time.Minute,
		4: 10 * time.Minute,
		5: 30 * time.Minute,
	}
)
