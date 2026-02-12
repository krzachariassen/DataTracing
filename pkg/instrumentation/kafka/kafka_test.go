package kafka

import (
	"context"
	"datatracing/pkg/sdk"
	"datatracing/pkg/tracing"
	"testing"
)

type producerSpy struct{ headers map[string]string }

func (p *producerSpy) Send(_ context.Context, _ string, _, _ []byte, headers map[string]string) error {
	p.headers = headers
	return nil
}

type noOpExporter struct{}

func (noOpExporter) Export(context.Context, tracing.Span) error { return nil }

func TestProducerInjectsHeaders(t *testing.T) {
	spy := &producerSpy{}
	tracer := sdk.NewTracer(noOpExporter{}, alwaysSampleSampler{})
	w := ProducerWrapper{Tracer: tracer, Producer: spy}
	if err := w.Send(context.Background(), "topic", nil, nil); err != nil {
		t.Fatal(err)
	}
	if spy.headers["trace-id"] == "" {
		t.Fatal("trace-id missing")
	}
}

func TestWrapConsumerExtractsContext(t *testing.T) {
	tracer := sdk.NewTracer(noOpExporter{}, alwaysSampleSampler{})
	called := false
	h := WrapConsumer(tracer, "topic", func(ctx context.Context, _, _ []byte, _ map[string]string) error {
		called = true
		return nil
	})
	err := h(context.Background(), nil, nil, map[string]string{"trace-id": "t1", "span-id": "s1", "sampled": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("handler not called")
	}
}

type alwaysSampleSampler struct{}

func (alwaysSampleSampler) ShouldSample(map[string]string, float64) bool { return true }
