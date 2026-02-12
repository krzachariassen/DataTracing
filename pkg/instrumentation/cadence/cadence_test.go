package cadence

import (
	"context"
	"datatracing/pkg/propagation"
	"datatracing/pkg/sdk"
	"datatracing/pkg/tracing"
	"testing"
)

type noOpExporter struct{}

func (noOpExporter) Export(context.Context, tracing.Span) error { return nil }

func TestInjectExtractWorkflow(t *testing.T) {
	tracer := sdk.NewTracer(noOpExporter{}, alwaysSampleSampler{})
	ctx, span := tracer.Start(context.Background(), "wf", tracing.SpanKindTransform)
	defer span.End()

	wf := WorkflowContext{}
	InjectWorkflow(ctx, wf)
	ctx2 := ExtractWorkflow(context.Background(), wf)
	if tc, ok := propagation.FromContext(ctx2); !ok || tc.TraceID == "" {
		t.Fatal("expected trace context")
	}
}

func TestStartWorkflowAndActivitySpan(t *testing.T) {
	tracer := sdk.NewTracer(noOpExporter{}, alwaysSampleSampler{})
	ctx, wfSpan := StartWorkflowSpan(context.Background(), tracer, "wf1")
	wfSpan.End()
	_, actSpan := StartActivitySpan(ctx, tracer, "act1")
	actSpan.End()
}

type alwaysSampleSampler struct{}

func (alwaysSampleSampler) ShouldSample(map[string]string, float64) bool { return true }
