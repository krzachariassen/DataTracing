package application

import (
	"datatracing/internal/domain"
	"strings"
	"time"
)

type HeadSampler struct {
	Probability float64
	TenantKey   string
	Tenants     map[string]float64
	Entities    map[string]bool
}

func (s HeadSampler) ShouldSample(attributes map[string]string, randValue float64) bool {
	if len(s.Entities) > 0 {
		if entity, ok := attributes["entity_id"]; ok && s.Entities[entity] {
			return true
		}
	}
	if len(s.Tenants) > 0 && s.TenantKey != "" {
		if tenant, ok := attributes[s.TenantKey]; ok {
			if p, ok := s.Tenants[tenant]; ok {
				return randValue <= p
			}
		}
	}
	if s.Probability <= 0 {
		return false
	}
	if s.Probability >= 1 {
		return true
	}
	return randValue <= s.Probability
}

type TailSamplingPolicy struct {
	ErrorAlways      bool
	LatencyThreshold time.Duration
	AttributeMatch   map[string]string
	MaxSpanCount     int
}

func (p TailSamplingPolicy) Keep(spans []domain.Span) bool {
	if len(spans) == 0 {
		return false
	}
	if p.MaxSpanCount > 0 && len(spans) >= p.MaxSpanCount {
		return true
	}
	for _, span := range spans {
		if p.ErrorAlways && span.Status == domain.SpanStatusError {
			return true
		}
		if p.LatencyThreshold > 0 && span.EndTime.Sub(span.StartTime) >= p.LatencyThreshold {
			return true
		}
		if len(p.AttributeMatch) > 0 {
			for k, v := range p.AttributeMatch {
				if strings.EqualFold(span.Attributes[k], v) {
					return true
				}
			}
		}
	}
	for _, span := range spans {
		if span.Sampled {
			return true
		}
	}
	return false
}
