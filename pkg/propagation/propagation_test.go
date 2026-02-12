package propagation

import (
	"context"
	"testing"
)

type mapCarrier map[string]string

func (m mapCarrier) Get(key string) string { return m[key] }
func (m mapCarrier) Set(key, value string) { m[key] = value }

func TestInjectExtract(t *testing.T) {
	ctx := WithTraceContext(context.Background(), TraceContext{TraceID: "t", SpanID: "s", ParentID: "p", Sampled: true})
	c := mapCarrier{}
	Inject(ctx, c)
	ctx2 := Extract(context.Background(), c)
	tc, ok := FromContext(ctx2)
	if !ok || tc.TraceID != "t" || !tc.Sampled {
		t.Fatalf("unexpected tc: %+v ok=%v", tc, ok)
	}
}
