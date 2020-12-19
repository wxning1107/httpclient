package slidingwindow

import (
	"container/ring"
	"sync"
	"time"
)

type SlidingWindow struct {
	bucketRing     *ring.Ring
	bucketInterval time.Duration
	bucketRWMutex  sync.RWMutex
	lastAccessTime time.Time
}

func NewSlidingWindow(bucketInterval time.Duration, bucketCount int) *SlidingWindow {
	buckets := ring.New(bucketCount)
	for i := 0; i < buckets.Len(); i++ {
		buckets.Value = new(Bucket)
		buckets = buckets.Next()
	}

	return &SlidingWindow{
		bucketRing:     buckets,
		bucketInterval: bucketInterval,
		lastAccessTime: time.Now(),
	}
}

func (sw *SlidingWindow) getCurrentBucket() *Bucket {
	currentBucket := sw.bucketRing.Value.(*Bucket)
	timeDiff := time.Since(sw.lastAccessTime)
	if timeDiff > sw.bucketInterval {
		for i := 0; i < sw.bucketRing.Len(); i++ {
			sw.bucketRing = sw.bucketRing.Next()
			currentBucket = sw.bucketRing.Value.(*Bucket)
			currentBucket.reset()
			timeDiff = time.Duration(int64(timeDiff) - int64(sw.bucketInterval))
			if timeDiff < sw.bucketInterval {
				break
			}
		}
		sw.lastAccessTime = time.Now()
	}

	return currentBucket
}

func (sw *SlidingWindow) Fail() {
	sw.bucketRWMutex.Lock()
	defer sw.bucketRWMutex.Unlock()

	bucket := sw.getCurrentBucket()
	bucket.fail()
}

func (sw *SlidingWindow) Success() {
	sw.bucketRWMutex.Lock()
	defer sw.bucketRWMutex.Unlock()
	bucket := sw.getCurrentBucket()
	bucket.success()
}

func (sw *SlidingWindow) Reset() {
	sw.bucketRWMutex.Lock()
	defer sw.bucketRWMutex.Unlock()
	sw.bucketRing.Do(func(x interface{}) {
		x.(*Bucket).reset()
	})
}

func (sw *SlidingWindow) FailureCount() (failureCount int64) {
	sw.bucketRWMutex.RLock()
	defer sw.bucketRWMutex.RUnlock()

	sw.bucketRing.Do(func(x interface{}) {
		bucket := x.(*Bucket)
		failureCount += bucket.failureCount
	})

	return
}

func (sw *SlidingWindow) SuccessCount() (successCount int64) {
	sw.bucketRWMutex.RLock()
	defer sw.bucketRWMutex.RUnlock()

	sw.bucketRing.Do(func(x interface{}) {
		bucket := x.(*Bucket)
		successCount += bucket.successCount
	})

	return
}

func (sw *SlidingWindow) FailureRate() float64 {
	sw.bucketRWMutex.RLock()
	defer sw.bucketRWMutex.RUnlock()

	var totalCount int64
	var failureCount int64
	sw.bucketRing.Do(func(x interface{}) {
		bucket := x.(*Bucket)
		totalCount += bucket.failureCount + bucket.successCount
		failureCount += bucket.failureCount
	})

	if totalCount == 0.0 {
		return 0.0
	}

	return float64(failureCount) / float64(totalCount)
}
