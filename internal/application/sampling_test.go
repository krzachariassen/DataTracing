package application

import (
	"datatracing/internal/domain"
	"testing"
	"time"
)

func TestHeadSampler_Table(t *testing.T) {
	s := HeadSampler{Probability: 0.1, TenantKey: "tenant", Tenants: map[string]float64{"gold": 1.0}, Entities: map[string]bool{"e1": true}}
	tests := []struct {
		name  string
		attrs map[string]string
		rand  float64
		want  bool
	}{
		{"entity override", map[string]string{"entity_id": "e1"}, 0.99, true},
		{"tenant override", map[string]string{"tenant": "gold"}, 0.99, true},
		{"prob sample", map[string]string{}, 0.01, true},
		{"prob drop", map[string]string{}, 0.9, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.ShouldSample(tt.attrs, tt.rand); got != tt.want {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestTailSamplingPolicy_Keep(t *testing.T) {
	now := time.Now()
	policy := TailSamplingPolicy{ErrorAlways: true, LatencyThreshold: 100 * time.Millisecond, AttributeMatch: map[string]string{"priority": "high"}, MaxSpanCount: 5}
	if !policy.Keep([]domain.Span{{Status: domain.SpanStatusError}}) {
		t.Fatal("error span should be kept")
	}
	if !policy.Keep([]domain.Span{{StartTime: now, EndTime: now.Add(200 * time.Millisecond)}}) {
		t.Fatal("high latency should be kept")
	}
	if !policy.Keep([]domain.Span{{Attributes: map[string]string{"priority": "HIGH"}}}) {
		t.Fatal("attribute match should be kept")
	}
}
