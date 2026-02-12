package application

import (
	"context"
	"datatracing/internal/domain"
	"sort"
)

type QueryService struct {
	store TraceStore
}

func NewQueryService(store TraceStore) *QueryService {
	return &QueryService{store: store}
}

func (s *QueryService) GetTraceDAG(ctx context.Context, traceID string) (domain.TraceDAG, error) {
	spans, err := s.store.GetTrace(ctx, traceID)
	if err != nil {
		return domain.TraceDAG{}, err
	}
	return BuildDAG(spans), nil
}

func (s *QueryService) SearchTraces(ctx context.Context, filter domain.QueryFilter) ([]domain.TraceSummary, error) {
	return s.store.QueryTraces(ctx, filter)
}

func BuildDAG(spans []domain.Span) domain.TraceDAG {
	nodes := make(map[string]*domain.DAGNode, len(spans))
	allLinks := make([]domain.SpanLink, 0)
	traceID := ""

	sort.Slice(spans, func(i, j int) bool { return spans[i].StartTime.Before(spans[j].StartTime) })

	for _, span := range spans {
		s := span
		if traceID == "" {
			traceID = s.TraceID
		}
		nodes[s.SpanID] = &domain.DAGNode{Span: s}
		allLinks = append(allLinks, s.Links...)
	}

	roots := make([]*domain.DAGNode, 0)
	orphans := make([]*domain.DAGNode, 0)
	allNodes := make([]*domain.DAGNode, 0, len(nodes))

	for _, span := range spans {
		node := nodes[span.SpanID]
		allNodes = append(allNodes, node)
		if span.ParentID == "" {
			roots = append(roots, node)
			continue
		}
		parent, ok := nodes[span.ParentID]
		if !ok {
			orphans = append(orphans, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}

	return domain.TraceDAG{TraceID: traceID, Roots: roots, Nodes: allNodes, Links: allLinks, Orphans: orphans}
}
