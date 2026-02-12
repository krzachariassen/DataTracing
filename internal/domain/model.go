package domain

import "datatracing/pkg/tracing"

type SpanKind = tracing.SpanKind

type SpanStatus = tracing.SpanStatus

type LinkType = tracing.LinkType

type Trace = tracing.Trace

type SpanEvent = tracing.SpanEvent

type SpanLink = tracing.SpanLink

type Span = tracing.Span

type QueryFilter = tracing.QueryFilter

type TraceSummary = tracing.TraceSummary

type TraceDAG = tracing.TraceDAG

type DAGNode = tracing.DAGNode

const (
	SpanKindIngest      = tracing.SpanKindIngest
	SpanKindTransform   = tracing.SpanKindTransform
	SpanKindEnrich      = tracing.SpanKindEnrich
	SpanKindValidate    = tracing.SpanKindValidate
	SpanKindJoin        = tracing.SpanKindJoin
	SpanKindAggregate   = tracing.SpanKindAggregate
	SpanKindMaterialize = tracing.SpanKindMaterialize
	SpanKindPublish     = tracing.SpanKindPublish
	SpanKindServe       = tracing.SpanKindServe
)

const (
	SpanStatusOK        = tracing.SpanStatusOK
	SpanStatusError     = tracing.SpanStatusError
	SpanStatusTimeout   = tracing.SpanStatusTimeout
	SpanStatusCancelled = tracing.SpanStatusCancelled
)

const (
	LinkTypeChildOf        = tracing.LinkTypeChildOf
	LinkTypeFollowsFrom    = tracing.LinkTypeFollowsFrom
	LinkTypeDataDependency = tracing.LinkTypeDataDependency
)
