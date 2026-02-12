package domain

import "testing"

func TestConstants_AreStable(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"kind transform", string(SpanKindTransform), "TRANSFORM"},
		{"status error", string(SpanStatusError), "ERROR"},
		{"link follows", string(LinkTypeFollowsFrom), "FOLLOWS_FROM"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %q want %q", tt.got, tt.want)
			}
		})
	}
}
