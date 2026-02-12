package middleware

import (
	"datatracing/pkg/propagation"
	"datatracing/pkg/sdk"
	"datatracing/pkg/tracing"
	"net/http"
)

type HeaderCarrier struct{ H http.Header }

func (c HeaderCarrier) Get(key string) string { return c.H.Get(key) }
func (c HeaderCarrier) Set(key, value string) { c.H.Set(key, value) }

func TraceHTTP(tracer *sdk.Tracer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := propagation.Extract(r.Context(), HeaderCarrier{H: r.Header})
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path, tracing.SpanKindServe)
		defer span.End()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
