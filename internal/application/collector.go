package application

import (
	"context"
	"datatracing/internal/domain"
	"sync"
	"time"
)

type CollectorService struct {
	store      TraceStore
	policy     TailSamplingPolicy
	batchSize  int
	flushEvery time.Duration

	ch chan domain.Span
	wg sync.WaitGroup
}

func NewCollectorService(store TraceStore, policy TailSamplingPolicy, batchSize int, flushEvery time.Duration, workers int, buffer int) *CollectorService {
	if batchSize <= 0 {
		batchSize = 100
	}
	if flushEvery <= 0 {
		flushEvery = 250 * time.Millisecond
	}
	if workers <= 0 {
		workers = 2
	}
	if buffer <= 0 {
		buffer = 1000
	}
	c := &CollectorService{store: store, policy: policy, batchSize: batchSize, flushEvery: flushEvery, ch: make(chan domain.Span, buffer)}
	c.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go c.worker()
	}
	return c
}

func (c *CollectorService) Ingest(span domain.Span) bool {
	select {
	case c.ch <- span:
		return true
	default:
		return false
	}
}

func (c *CollectorService) Close() {
	close(c.ch)
	c.wg.Wait()
}

func (c *CollectorService) worker() {
	defer c.wg.Done()
	ticker := time.NewTicker(c.flushEvery)
	defer ticker.Stop()

	batch := make([]domain.Span, 0, c.batchSize)
	flush := func() {
		if len(batch) == 0 {
			return
		}
		byTrace := map[string][]domain.Span{}
		for _, span := range batch {
			byTrace[span.TraceID] = append(byTrace[span.TraceID], span)
		}
		for _, traceSpans := range byTrace {
			if c.policy.Keep(traceSpans) {
				_ = c.store.SaveSpans(context.Background(), traceSpans)
			}
		}
		batch = batch[:0]
	}

	for {
		select {
		case span, ok := <-c.ch:
			if !ok {
				flush()
				return
			}
			batch = append(batch, span)
			if len(batch) >= c.batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}
