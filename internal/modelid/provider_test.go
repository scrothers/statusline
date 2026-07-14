package modelid

import "testing"

func TestDetectProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id       string
		wantProv Provider
		wantOK   bool
	}{
		{"bedrock plain", "anthropic.claude-opus-4-8", ProviderAWS, true},
		{
			"bedrock cross-region dated",
			"us.anthropic.claude-3-5-sonnet-20241022-v2:0",
			ProviderAWS,
			true,
		},
		{
			"bedrock full arn",
			"arn:aws:bedrock:us-east-1::foundation-model/anthropic.claude-3-5-sonnet-20241022-v2:0",
			ProviderAWS,
			true,
		},
		{"vertex dated snapshot", "claude-opus-4-5@20251101", ProviderGCP, true},
		{
			"vertex full resource path",
			"projects/my-proj/locations/us-central1/publishers/anthropic/models/claude-opus-4-5@20251101",
			ProviderGCP,
			true,
		},
		{"openrouter style", "anthropic/claude-3.5-sonnet:beta", ProviderRouter, true},
		{"bare first-party id", "claude-opus-4-8", "", false},
		{"gateway reordered id, still bare", "claude-4-8-opus", "", false},
		{"empty id", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotProv, gotOK := DetectProvider(tt.id)
			if gotOK != tt.wantOK {
				t.Fatalf("DetectProvider(%q) ok = %v, want %v", tt.id, gotOK, tt.wantOK)
			}
			if gotProv != tt.wantProv {
				t.Errorf("DetectProvider(%q) = %q, want %q", tt.id, gotProv, tt.wantProv)
			}
		})
	}
}

func TestParseProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       string
		wantProv Provider
		wantOK   bool
	}{
		{"aws", "aws", ProviderAWS, true},
		{"gcp", "gcp", ProviderGCP, true},
		{"azure", "azure", ProviderAzure, true},
		{"router", "router", ProviderRouter, true},
		{"gateway", "gateway", ProviderGateway, true},
		{"uppercase", "AWS", ProviderAWS, true},
		{"surrounding whitespace", "  gcp  ", ProviderGCP, true},
		{"empty", "", "", false},
		{"unrecognized", "not-a-provider", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotProv, gotOK := ParseProvider(tt.in)
			if gotOK != tt.wantOK {
				t.Fatalf("ParseProvider(%q) ok = %v, want %v", tt.in, gotOK, tt.wantOK)
			}
			if gotProv != tt.wantProv {
				t.Errorf("ParseProvider(%q) = %q, want %q", tt.in, gotProv, tt.wantProv)
			}
		})
	}
}
