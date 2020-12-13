package circuitbreaker

type TripFunc func(*Breaker) bool

func RateTripFunc(rate float64, minSamples int64) TripFunc {
	return func(b *Breaker) bool {
		samples := b.FailureCount() + b.SuccessCount()
		return samples >= minSamples && b.FailureRate() >= rate
	}
}
