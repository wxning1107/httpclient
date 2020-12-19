package circuitbreaker

import (
	"context"
	"github.com/cenkalti/backoff"
	"httpclient/slidingwindow"
	"sync/atomic"
	"time"
)

type BreakerStatus = int

const (
	BreakerNotReady BreakerStatus = iota + 1
	BreakerTimeout
	BreakerSuccess
	BreakerFailed
)

type BreakerState = int64

const (
	ClosedState BreakerState = iota + 1
	HalfOpenState
	OpenState
)

type Breaker struct {
	backOff     backoff.BackOff
	nextBackOff time.Duration

	ShouldTrip   TripFunc
	ShouldResume ShouldResume
	counts       *slidingwindow.SlidingWindow

	breakerState             int64
	lastRequestTime          int64
	attemptHalfOpensInterval time.Duration
	halfOpenRequesting       int64
	continuousFailuresCount  int64
	continuousSuccessCount   int64
}

func (b *Breaker) trip() {
	atomic.StoreInt64(&b.breakerState, OpenState)
	b.counts.Reset()
}

func (b *Breaker) resume() {
	atomic.StoreInt64(&b.breakerState, ClosedState)
	b.counts.Reset()
}

func (b *Breaker) Call(ctx context.Context, f func() error, timeout time.Duration) (state BreakerStatus, err error) {
	if !b.Ready() {
		return BreakerNotReady, nil
	}

	if timeout == 0 {
		err = f()
	} else {
		errCh := make(chan error, 1)
		go func() {
			errCh <- f()
			close(errCh)
		}()

		select {
		case err = <-errCh:
		case <-time.After(timeout):
			return BreakerTimeout, nil
		}
	}

	if err != nil {
		if ctx.Err() == context.Canceled {
			b.fail()
		}
		return BreakerFailed, err
	}
	b.success()

	return BreakerSuccess, nil
}

func (b *Breaker) Ready() bool {
	currentState := b.currentState()
	switch currentState {
	case ClosedState:
		return true
	case OpenState:
		return false
	case HalfOpenState:
		if time.Since(time.Unix(0, atomic.LoadInt64(&b.lastRequestTime))) >= b.nextBackOff &&
			atomic.CompareAndSwapInt64(&b.halfOpenRequesting, 0, 1) {
			b.nextBackOff = b.backOff.NextBackOff()
			return true
		}
		return false
	}

	return false
}

func (b *Breaker) currentState() BreakerState {
	breakerState := atomic.LoadInt64(&b.breakerState)
	switch breakerState {
	case ClosedState:
		return ClosedState
	case HalfOpenState:
		return HalfOpenState
	case OpenState:
		if time.Since(time.Unix(0, atomic.LoadInt64(&b.lastRequestTime))) > b.attemptHalfOpensInterval &&
			atomic.CompareAndSwapInt64(&b.breakerState, OpenState, HalfOpenState) {
			return HalfOpenState
		}
		return OpenState
	}

	return OpenState
}

func (b *Breaker) success() {
	b.counts.Success()
	atomic.AddInt64(&b.continuousSuccessCount, 1)
	atomic.StoreInt64(&b.continuousFailuresCount, 0)

	switch b.currentState() {
	case ClosedState:
		atomic.StoreInt64(&b.lastRequestTime, time.Now().UnixNano())
		return
	case HalfOpenState:
		atomic.StoreInt64(&b.halfOpenRequesting, 0)
		if b.ShouldResume != nil && b.ShouldResume(b) {
			b.resume()
		}
	}

	b.backOff.Reset()
	b.nextBackOff = b.backOff.NextBackOff()
	atomic.StoreInt64(&b.lastRequestTime, time.Now().UnixNano())
}

func (b *Breaker) fail() {
	b.counts.Fail()
	atomic.AddInt64(&b.continuousFailuresCount, 1)
	atomic.StoreInt64(&b.continuousSuccessCount, 0)

	switch b.currentState() {
	case ClosedState:
		if b.ShouldTrip != nil && b.ShouldTrip(b) {
			b.trip()
		}
	case HalfOpenState:
		atomic.StoreInt64(&b.halfOpenRequesting, 0)
		b.trip()
	}

	atomic.StoreInt64(&b.lastRequestTime, time.Now().UnixNano())
}

func (b *Breaker) FailureCount() int64 {
	return b.counts.FailureCount()
}

func (b *Breaker) SuccessCount() int64 {
	return b.counts.SuccessCount()
}

func (b *Breaker) FailureRate() float64 {
	return b.counts.FailureRate()
}

func (b *Breaker) ContinuousSuccessCount() int64 {
	return atomic.LoadInt64(&b.continuousSuccessCount)
}
