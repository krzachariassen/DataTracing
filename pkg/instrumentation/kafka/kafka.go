package kafka

import (
	"context"
	"datatracing/pkg/propagation"
	"datatracing/pkg/sdk"
	"datatracing/pkg/tracing"
)

type HeaderCarrier map[string]string

func (h HeaderCarrier) Get(key string) string { return h[key] }
func (h HeaderCarrier) Set(key, value string) { h[key] = value }

type Producer interface {
	Send(ctx context.Context, topic string, key, value []byte, headers map[string]string) error
}

type ConsumerHandler func(ctx context.Context, key, value []byte, headers map[string]string) error

type ProducerWrapper struct {
	Tracer   *sdk.Tracer
	Producer Producer
}

func (w ProducerWrapper) Send(ctx context.Context, topic string, key, value []byte) error {
	ctx, span := w.Tracer.Start(ctx, "kafka.produce:"+topic, tracing.SpanKindPublish)
	defer span.End()
	headers := map[string]string{}
	propagation.Inject(ctx, HeaderCarrier(headers))
	return w.Producer.Send(ctx, topic, key, value, headers)
}

func WrapConsumer(tracer *sdk.Tracer, topic string, next ConsumerHandler) ConsumerHandler {
	return func(ctx context.Context, key, value []byte, headers map[string]string) error {
		ctx = propagation.Extract(ctx, HeaderCarrier(headers))
		ctx, span := tracer.Start(ctx, "kafka.consume:"+topic, tracing.SpanKindIngest)
		defer span.End()
		return next(ctx, key, value, headers)
	}
}
