package sdk

import (
	"context"
	"datatracing/pkg/propagation"
	"datatracing/pkg/tracing"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Exporter interface {
	Export(ctx context.Context, span tracing.Span) error
}

type Sampler interface {
	ShouldSample(attributes map[string]string, randValue float64) bool
}

type alwaysOnSampler struct{}

func (alwaysOnSampler) ShouldSample(map[string]string, float64) bool { return true }

type Tracer struct {
	exporter    Exporter
	headSampler Sampler
	counter     atomic.Uint64
}

func NewTracer(exporter Exporter, sampler Sampler) *Tracer {
	if sampler == nil {
		sampler = alwaysOnSampler{}
	}
	return &Tracer{exporter: exporter, headSampler: sampler}
}

type Span struct {
	mu       sync.Mutex
	ended    bool
	tracer   *Tracer
	ctx      context.Context
	span     tracing.Span
	parentTC propagation.TraceContext
}

func (t *Tracer) Start(ctx context.Context, operation string, kind tracing.SpanKind) (context.Context, *Span) {
	now := time.Now().UTC()
	parentTC, _ := propagation.FromContext(ctx)
	traceID := parentTC.TraceID
	if traceID == "" {
		traceID = t.newID("trace")
	}
	spanID := t.newID("span")
	attrs := map[string]string{}
	sampled := t.headSampler.ShouldSample(attrs, rand.Float64())
	if parentTC.TraceID != "" {
		sampled = parentTC.Sampled
	}
	ds := tracing.Span{
		TraceID:    traceID,
		SpanID:     spanID,
		ParentID:   parentTC.SpanID,
		Operation:  operation,
		Kind:       kind,
		StartTime:  now,
		Status:     tracing.SpanStatusOK,
		Attributes: attrs,
		Sampled:    sampled,
	}
	tc := propagation.TraceContext{TraceID: traceID, SpanID: spanID, ParentID: parentTC.SpanID, Sampled: sampled}
	childCtx := propagation.WithTraceContext(ctx, tc)
	return childCtx, &Span{tracer: t, ctx: childCtx, span: ds, parentTC: parentTC}
}

func (s *Span) End() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ended {
		return
	}
	s.ended = true
	s.span.EndTime = time.Now().UTC()
	if s.span.EndTime.Before(s.span.StartTime) {
		s.span.EndTime = s.span.StartTime
	}
	if s.tracer.exporter != nil {
		_ = s.tracer.exporter.Export(s.ctx, s.span)
	}
}

func (s *Span) SetAttribute(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.span.Attributes[key] = value
}

func (s *Span) AddEvent(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.span.Events = append(s.span.Events, tracing.SpanEvent{Name: name, Timestamp: time.Now().UTC()})
}

func (s *Span) SetStatus(status tracing.SpanStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.span.Status = status
}

func (t *Tracer) newID(prefix string) string {
	n := t.counter.Add(1)
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UTC().UnixNano(), n)
}
