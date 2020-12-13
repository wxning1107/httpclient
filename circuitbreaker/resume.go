package circuitbreaker

import (
	"sync/atomic"
	"time"
)

type ShouldResume func(*Breaker) bool

func ContinuousOrIntervalResumeFunc(interval time.Duration, threshold int64) ShouldResume {
	return func(b *Breaker) bool {
		return b.ContinuousSuccessCount() >= threshold ||
			time.Now().Sub(time.Unix(0, atomic.LoadInt64(&b.lastRequestTime))) >= interval
	}
}
