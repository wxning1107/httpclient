package circuitbreaker

import (
	"github.com/cenkalti/backoff"
	"httpclient/slidingwindow"
	"time"
)

const (
	DefaultExponentialBackOffInterval       = time.Millisecond * 500
	DefaultExponentialBackOffMaxElapsedTime = 0

	DefaultContinuousSuccessResumeThreshold = 3
	DefaultContinuousSuccessResumeInterval  = time.Second * 3

	DefaultWindowSlideInterval = time.Second
	DefaultWindowBucketCount   = 10

	DefaultAttemptHalfOpensInterval = time.Second
)

type Option struct {
	BackOff                  backoff.BackOff
	ShouldTrip               TripFunc
	ShouldResume             ShouldResume
	windowSlideInterval      time.Duration
	windowBucketCount        int
	attemptHalfOpensInterval time.Duration
}

func NewBreakerWithOptions(options *Option) *Breaker {
	if options == nil {
		options = &Option{}
	}

	if options.BackOff == nil {
		options.BackOff = NewExponentialBackOff(DefaultExponentialBackOffInterval, DefaultExponentialBackOffMaxElapsedTime)
	}

	if options.ShouldResume == nil {
		options.ShouldResume = ContinuousOrIntervalResumeFunc(DefaultContinuousSuccessResumeInterval, DefaultContinuousSuccessResumeThreshold)
	}

	if options.windowSlideInterval == 0 {
		options.windowSlideInterval = DefaultWindowSlideInterval
	}

	if options.attemptHalfOpensInterval == 0 {
		options.attemptHalfOpensInterval = DefaultAttemptHalfOpensInterval
	}

	if options.windowBucketCount == 0 {
		options.windowBucketCount = DefaultWindowBucketCount
	}

	return &Breaker{
		backOff:                  options.BackOff,
		nextBackOff:              options.BackOff.NextBackOff(),
		ShouldTrip:               options.ShouldTrip,
		ShouldResume:             options.ShouldResume,
		lastRequestTime:          0,
		attemptHalfOpensInterval: options.attemptHalfOpensInterval,
		counts:                   slidingwindow.NewSlidingWindow(options.windowSlideInterval, options.windowBucketCount),
	}

}

func NewRateBreaker(rate float64, minSamples int64) *Breaker {
	return NewBreakerWithOptions(&Option{
		ShouldTrip: RateTripFunc(rate, minSamples),
	})
}

func NewExponentialBackOff(interval, maxElapsedTime time.Duration) backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = interval
	b.MaxElapsedTime = maxElapsedTime
	b.Reset()
	return b
}
