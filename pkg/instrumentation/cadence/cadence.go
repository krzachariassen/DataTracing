package cadence

import (
	"context"
	"datatracing/pkg/propagation"
	"datatracing/pkg/sdk"
	"datatracing/pkg/tracing"
)

type WorkflowContext map[string]string

func (c WorkflowContext) Get(key string) string { return c[key] }
func (c WorkflowContext) Set(key, value string) { c[key] = value }

func InjectWorkflow(ctx context.Context, wfCtx WorkflowContext) {
	propagation.Inject(ctx, wfCtx)
}

func ExtractWorkflow(ctx context.Context, wfCtx WorkflowContext) context.Context {
	return propagation.Extract(ctx, wfCtx)
}

func StartWorkflowSpan(ctx context.Context, tracer *sdk.Tracer, workflowName string) (context.Context, *sdk.Span) {
	return tracer.Start(ctx, "cadence.workflow:"+workflowName, tracing.SpanKindTransform)
}

func StartActivitySpan(ctx context.Context, tracer *sdk.Tracer, activityName string) (context.Context, *sdk.Span) {
	return tracer.Start(ctx, "cadence.activity:"+activityName, tracing.SpanKindEnrich)
}
