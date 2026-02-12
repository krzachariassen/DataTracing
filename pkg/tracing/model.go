package tracing

import "time"

type SpanKind string

const (
	SpanKindIngest      SpanKind = "INGEST"
	SpanKindTransform   SpanKind = "TRANSFORM"
	SpanKindEnrich      SpanKind = "ENRICH"
	SpanKindValidate    SpanKind = "VALIDATE"
	SpanKindJoin        SpanKind = "JOIN"
	SpanKindAggregate   SpanKind = "AGGREGATE"
	SpanKindMaterialize SpanKind = "MATERIALIZE"
	SpanKindPublish     SpanKind = "PUBLISH"
	SpanKindServe       SpanKind = "SERVE"
)

type SpanStatus string

const (
	SpanStatusOK        SpanStatus = "OK"
	SpanStatusError     SpanStatus = "ERROR"
	SpanStatusTimeout   SpanStatus = "TIMEOUT"
	SpanStatusCancelled SpanStatus = "CANCELLED"
)

type LinkType string

const (
	LinkTypeChildOf        LinkType = "CHILD_OF"
	LinkTypeFollowsFrom    LinkType = "FOLLOWS_FROM"
	LinkTypeDataDependency LinkType = "DATA_DEPENDENCY"
)

type Trace struct {
	TraceID    string    `json:"trace_id"`
	RootSpanID string    `json:"root_span_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

type SpanEvent struct {
	Name       string            `json:"name"`
	Timestamp  time.Time         `json:"timestamp"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type SpanLink struct {
	TraceID string   `json:"trace_id"`
	SpanID  string   `json:"span_id"`
	Type    LinkType `json:"type"`
}

type Span struct {
	TraceID    string            `json:"trace_id"`
	SpanID     string            `json:"span_id"`
	ParentID   string            `json:"parent_id,omitempty"`
	Operation  string            `json:"operation"`
	Kind       SpanKind          `json:"kind"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
	Status     SpanStatus        `json:"status"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Events     []SpanEvent       `json:"events,omitempty"`
	Links      []SpanLink        `json:"links,omitempty"`
	Sampled    bool              `json:"sampled"`
}

type QueryFilter struct {
	Operation string
	Status    SpanStatus
	From      time.Time
	To        time.Time
	Limit     int
}

type TraceSummary struct {
	TraceID   string    `json:"trace_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	SpanCount int       `json:"span_count"`
	Status    string    `json:"status"`
}

type TraceDAG struct {
	TraceID string     `json:"trace_id"`
	Roots   []*DAGNode `json:"roots"`
	Nodes   []*DAGNode `json:"nodes"`
	Links   []SpanLink `json:"links"`
	Orphans []*DAGNode `json:"orphans,omitempty"`
}

type DAGNode struct {
	Span     Span       `json:"span"`
	Children []*DAGNode `json:"children,omitempty"`
}
