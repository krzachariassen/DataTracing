package sdk

import (
	"bytes"
	"context"
	"datatracing/pkg/tracing"
	"encoding/json"
	"net/http"
)

type HTTPExporter struct {
	Client *http.Client
	URL    string
}

func (e HTTPExporter) Export(ctx context.Context, span tracing.Span) error {
	b, err := json.Marshal(span)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := e.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
