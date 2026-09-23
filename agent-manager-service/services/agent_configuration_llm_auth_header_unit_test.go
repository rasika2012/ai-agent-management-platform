// Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package services

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wso2/agent-manager/agent-manager-service/models"
)

// The header an agent sends its credential in is reported from the proxy's own
// stored security config rather than a literal, so a proxy provisioned with a
// non-default header is described accurately in the config response.
func TestLLMProxyAPIKeyHeaderName(t *testing.T) {
	proxyWithHeader := func(key string) *models.LLMProxy {
		return &models.LLMProxy{
			Configuration: models.LLMProxyConfig{
				Security: &models.SecurityConfig{
					APIKey: &models.APIKeySecurity{Key: key},
				},
			},
		}
	}

	tests := []struct {
		name  string
		proxy *models.LLMProxy
		want  string
	}{
		{
			name:  "reports the stored header",
			proxy: proxyWithHeader("X-Custom-Key"),
			want:  "X-Custom-Key",
		},
		{
			// A proxy's deployed gateway policy enforces whatever header it was
			// stored with, so the stored value has to win over the default even
			// when the default later changes.
			name:  "stored header wins over the default",
			proxy: proxyWithHeader("X-API-Key"),
			want:  "X-API-Key",
		},
		{
			name:  "trims surrounding whitespace",
			proxy: proxyWithHeader("  x-api-key  "),
			want:  "x-api-key",
		},
		{
			name:  "falls back when the stored header is blank",
			proxy: proxyWithHeader("   "),
			want:  models.DefaultLLMProxyAPIKeyHeader,
		},
		{
			name:  "falls back when no api key security is configured",
			proxy: &models.LLMProxy{Configuration: models.LLMProxyConfig{Security: &models.SecurityConfig{}}},
			want:  models.DefaultLLMProxyAPIKeyHeader,
		},
		{
			name:  "falls back when the proxy has no security block",
			proxy: &models.LLMProxy{},
			want:  models.DefaultLLMProxyAPIKeyHeader,
		},
		{
			name:  "falls back for a nil proxy",
			proxy: nil,
			want:  models.DefaultLLMProxyAPIKeyHeader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, llmProxyAPIKeyHeaderName(tt.proxy))
		})
	}
}

// A proxy fronting a provider takes the provider's own configured header, so the
// name an admin sets on the provider is the one their agents authenticate with.
func TestProviderProxyAPIKeyHeader(t *testing.T) {
	providerWithAPIKey := func(apiKey *models.APIKeySecurity) *models.LLMProvider {
		return &models.LLMProvider{
			Configuration: models.LLMProviderConfig{
				Security: &models.SecurityConfig{APIKey: apiKey},
			},
		}
	}

	tests := []struct {
		name     string
		provider *models.LLMProvider
		want     string
	}{
		{
			name:     "follows the header configured on the provider",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "x-api-key", In: "header"}),
			want:     "x-api-key",
		},
		{
			name:     "treats an unset location as a header",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "x-api-key"}),
			want:     "x-api-key",
		},
		{
			// The proxy always takes its credential in a header, so a provider
			// that reads one from the query string names nothing usable here.
			name:     "falls back when the provider reads its key from the query",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "apikey", In: "query"}),
			want:     models.DefaultLLMProxyAPIKeyHeader,
		},
		{
			name:     "falls back when the provider names no header",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "  "}),
			want:     models.DefaultLLMProxyAPIKeyHeader,
		},
		{
			name:     "falls back when the provider has no api key security",
			provider: providerWithAPIKey(nil),
			want:     models.DefaultLLMProxyAPIKeyHeader,
		},
		{
			name:     "falls back for a provider with no security block",
			provider: &models.LLMProvider{},
			want:     models.DefaultLLMProxyAPIKeyHeader,
		},
		{
			name:     "falls back for a nil provider",
			provider: nil,
			want:     models.DefaultLLMProxyAPIKeyHeader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, providerProxyAPIKeyHeader(tt.provider))
		})
	}
}

// The staleness check decides whether a provider edit redeploys a proxy at all, so it
// has to be exact in both directions: miss a difference and agents keep sending a header
// the gateway no longer accepts; report one spuriously and an unrelated provider edit
// redeploys the whole fleet.
func TestProxyAuthHeadersStale(t *testing.T) {
	proxy := func(ingress string, upstream *string) *models.LLMProxy {
		config := models.LLMProxyConfig{}
		if ingress != "" {
			config.Security = &models.SecurityConfig{APIKey: &models.APIKeySecurity{Key: ingress}}
		}
		if upstream != nil {
			config.UpstreamAuth = &models.UpstreamAuth{Header: upstream}
		}
		return &models.LLMProxy{Configuration: config}
	}
	header := func(s string) *string { return &s }

	tests := []struct {
		name                      string
		proxy                     *models.LLMProxy
		ingress, upstream         string
		wantIngress, wantUpstream bool
	}{
		{
			name:    "both already match",
			proxy:   proxy("x-api-key", header("x-api-key")),
			ingress: "x-api-key", upstream: "x-api-key",
		},
		{
			name:    "ingress differs",
			proxy:   proxy("API-Key", header("x-api-key")),
			ingress: "x-api-key", upstream: "x-api-key",
			wantIngress: true,
		},
		{
			name:    "upstream differs",
			proxy:   proxy("x-api-key", header("API-Key")),
			ingress: "x-api-key", upstream: "x-api-key",
			wantUpstream: true,
		},
		{
			// The pre-existing gap: a proxy provisioned before an edit carries both
			// stale names, and fixing only one leaves the other hop broken.
			name:    "both differ",
			proxy:   proxy("API-Key", header("API-Key")),
			ingress: "x-api-key", upstream: "x-api-key",
			wantIngress: true, wantUpstream: true,
		},
		{
			name:    "upstream auth present but unset header",
			proxy:   proxy("x-api-key", nil),
			ingress: "x-api-key", upstream: "x-api-key",
		},
		{
			name:    "proxy without api key security is left alone",
			proxy:   proxy("", header("x-api-key")),
			ingress: "x-api-key", upstream: "x-api-key",
		},
		{
			// A provider with no api-key header of its own names nothing to forward,
			// so the upstream hop must not be rewritten to an empty header.
			name:    "provider names no upstream header",
			proxy:   proxy("x-api-key", header("API-Key")),
			ingress: "x-api-key", upstream: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIngress, gotUpstream := proxyAuthHeadersStale(tt.proxy, tt.ingress, tt.upstream)
			require.Equal(t, tt.wantIngress, gotIngress, "ingress staleness")
			require.Equal(t, tt.wantUpstream, gotUpstream, "upstream staleness")
		})
	}
}

// Provisioning and reporting must agree on the header name: a freshly built proxy
// config is what a later config response reads back, so the two cannot drift.
func TestLLMProxyProvisionedHeaderMatchesReportedHeader(t *testing.T) {
	enabled := true
	provisioned := &models.LLMProxy{
		Configuration: models.LLMProxyConfig{
			Security: &models.SecurityConfig{
				Enabled: &enabled,
				APIKey: &models.APIKeySecurity{
					Enabled: &enabled,
					Key:     models.DefaultLLMProxyAPIKeyHeader,
					In:      "header",
				},
			},
		},
	}

	require.Equal(t, models.DefaultLLMProxyAPIKeyHeader, llmProxyAPIKeyHeaderName(provisioned))
}
