package slidingwindow

type Bucket struct {
	failureCount int64
	successCount int64
}

func (b *Bucket) reset() {
	b.failureCount = 0
	b.successCount = 0
}

func (b *Bucket) fail() {
	b.failureCount++
}

func (b *Bucket) success() {
	b.successCount++
}
