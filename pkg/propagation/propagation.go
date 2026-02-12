package propagation

import "context"

type Carrier interface {
	Get(key string) string
	Set(key, value string)
}

type TraceContext struct {
	TraceID  string
	SpanID   string
	ParentID string
	Sampled  bool
}

type ctxKey struct{}

func WithTraceContext(ctx context.Context, tc TraceContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, tc)
}

func FromContext(ctx context.Context) (TraceContext, bool) {
	tc, ok := ctx.Value(ctxKey{}).(TraceContext)
	return tc, ok
}

func Inject(ctx context.Context, carrier Carrier) {
	tc, ok := FromContext(ctx)
	if !ok {
		return
	}
	carrier.Set("trace-id", tc.TraceID)
	carrier.Set("span-id", tc.SpanID)
	carrier.Set("parent-id", tc.ParentID)
	if tc.Sampled {
		carrier.Set("sampled", "1")
	} else {
		carrier.Set("sampled", "0")
	}
}

func Extract(ctx context.Context, carrier Carrier) context.Context {
	sampled := carrier.Get("sampled") == "1" || carrier.Get("sampled") == "true"
	tc := TraceContext{
		TraceID:  carrier.Get("trace-id"),
		SpanID:   carrier.Get("span-id"),
		ParentID: carrier.Get("parent-id"),
		Sampled:  sampled,
	}
	if tc.TraceID == "" {
		return ctx
	}
	return WithTraceContext(ctx, tc)
}
