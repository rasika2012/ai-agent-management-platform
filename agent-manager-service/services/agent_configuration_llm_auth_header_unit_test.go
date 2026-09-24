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

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/wso2/agent-manager/agent-manager-service/models"
	"github.com/wso2/agent-manager/agent-manager-service/repositories/repomocks"
)

// The name and location an agent must send its credential in are reported from the
// proxy's own stored security config rather than a literal, so a proxy provisioned
// with a non-default header — or configured to read its key from the query string —
// is described accurately in the config response. A proxy that doesn't actually
// require a credential must report neither, or the response fabricates an auth
// requirement that isn't real.
func TestLLMProxyAPIKeySecurity(t *testing.T) {
	trueVal := true
	falseVal := false
	proxyWith := func(key, in string, enabled *bool) *models.LLMProxy {
		return &models.LLMProxy{
			Configuration: models.LLMProxyConfig{
				Security: &models.SecurityConfig{
					Enabled: &trueVal,
					APIKey:  &models.APIKeySecurity{Enabled: enabled, Key: key, In: in},
				},
			},
		}
	}

	tests := []struct {
		name     string
		proxy    *models.LLMProxy
		wantName string
		wantIn   string
	}{
		{
			name:     "reports the stored header",
			proxy:    proxyWith("X-Custom-Key", "header", &trueVal),
			wantName: "X-Custom-Key",
			wantIn:   "header",
		},
		{
			// A proxy's deployed gateway policy enforces whatever header it was
			// stored with, so the stored value has to win over the default even
			// when the default later changes.
			name:     "stored header wins over the default",
			proxy:    proxyWith("X-API-Key", "header", &trueVal),
			wantName: "X-API-Key",
			wantIn:   "header",
		},
		{
			name:     "trims surrounding whitespace",
			proxy:    proxyWith("  x-api-key  ", "header", &trueVal),
			wantName: "x-api-key",
			wantIn:   "header",
		},
		{
			name:     "falls back when the stored header is blank",
			proxy:    proxyWith("   ", "header", &trueVal),
			wantName: models.DefaultLLMProxyAPIKeyHeader,
			wantIn:   "header",
		},
		{
			// A query-based proxy's real parameter name must still be reported —
			// not silently replaced by a header-oriented default (CRIT: this used
			// to be mis-described as "header" with a fabricated name).
			name:     "reports a query-based key with its real name and location",
			proxy:    proxyWith("apikey", "query", &trueVal),
			wantName: "apikey",
			wantIn:   "query",
		},
		{
			name:     "falls back to the header default when the query key is blank",
			proxy:    proxyWith("", "query", &trueVal),
			wantName: models.DefaultLLMProxyAPIKeyHeader,
			wantIn:   "query",
		},
		{
			name:     "treats an unrecognized location as a header",
			proxy:    proxyWith("X-Custom-Key", "cookie", &trueVal),
			wantName: "X-Custom-Key",
			wantIn:   "header",
		},
		{
			// CRIT: this used to still fabricate a default "API-Key"/header pair
			// for a proxy that explicitly turned api-key auth off, incorrectly
			// telling agents a credential was required when the proxy accepts
			// unauthenticated calls.
			name:     "reports nothing when api-key auth is explicitly disabled",
			proxy:    proxyWith("X-Custom-Key", "header", &falseVal),
			wantName: "",
			wantIn:   "",
		},
		{
			name:     "reports nothing when api-key auth is not opted into",
			proxy:    proxyWith("X-Custom-Key", "header", nil),
			wantName: "",
			wantIn:   "",
		},
		{
			name: "reports nothing when the proxy's security block is turned off",
			proxy: &models.LLMProxy{
				Configuration: models.LLMProxyConfig{
					Security: &models.SecurityConfig{
						Enabled: &falseVal,
						APIKey:  &models.APIKeySecurity{Enabled: &trueVal, Key: "X-Custom-Key", In: "header"},
					},
				},
			},
			wantName: "",
			wantIn:   "",
		},
		{
			name:     "reports nothing when no api key security is configured",
			proxy:    &models.LLMProxy{Configuration: models.LLMProxyConfig{Security: &models.SecurityConfig{}}},
			wantName: "",
			wantIn:   "",
		},
		{
			name:     "reports nothing when the proxy has no security block",
			proxy:    &models.LLMProxy{},
			wantName: "",
			wantIn:   "",
		},
		{
			name:     "reports nothing for a nil proxy",
			proxy:    nil,
			wantName: "",
			wantIn:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotIn := llmProxyAPIKeySecurity(tt.proxy)
			require.Equal(t, tt.wantName, gotName)
			require.Equal(t, tt.wantIn, gotIn)
		})
	}
}

// A proxy fronting a provider takes the provider's own configured name and location,
// so what an admin sets on the provider is what their agents authenticate with —
// including when the provider reads its own key from the query string, which used to
// be silently collapsed to a header-based default (CRIT: this was the actual cause of
// a proxy still reporting "header"/the default name after its provider was switched
// to query — ProvisionProxy and the resync path both went through this function).
func TestProviderProxyAPIKeySecurity(t *testing.T) {
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
		wantName string
		wantIn   string
	}{
		{
			name:     "follows the header configured on the provider",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "x-api-key", In: "header"}),
			wantName: "x-api-key",
			wantIn:   "header",
		},
		{
			name:     "treats an unset location as a header",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "x-api-key"}),
			wantName: "x-api-key",
			wantIn:   "header",
		},
		{
			name:     "follows the provider's real name and location when it reads its key from the query",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "apikey", In: "query"}),
			wantName: "apikey",
			wantIn:   "query",
		},
		{
			name:     "falls back to the header default when the query key is blank",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "  ", In: "query"}),
			wantName: models.DefaultLLMProxyAPIKeyHeader,
			wantIn:   "query",
		},
		{
			name:     "falls back when the provider names no header",
			provider: providerWithAPIKey(&models.APIKeySecurity{Key: "  "}),
			wantName: models.DefaultLLMProxyAPIKeyHeader,
			wantIn:   "header",
		},
		{
			name:     "falls back when the provider has no api key security",
			provider: providerWithAPIKey(nil),
			wantName: models.DefaultLLMProxyAPIKeyHeader,
			wantIn:   "header",
		},
		{
			name:     "falls back for a provider with no security block",
			provider: &models.LLMProvider{},
			wantName: models.DefaultLLMProxyAPIKeyHeader,
			wantIn:   "header",
		},
		{
			name:     "falls back for a nil provider",
			provider: nil,
			wantName: models.DefaultLLMProxyAPIKeyHeader,
			wantIn:   "header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotIn := providerProxyAPIKeySecurity(tt.provider)
			require.Equal(t, tt.wantName, gotName)
			require.Equal(t, tt.wantIn, gotIn)
		})
	}
}

// The staleness check decides whether a provider edit redeploys a proxy at all, so it
// has to be exact in both directions: miss a difference and agents keep sending a header
// the gateway no longer accepts; report one spuriously and an unrelated provider edit
// redeploys the whole fleet. Ingress staleness must catch a location-only change too —
// a provider switched from header to query with the same key name still needs its
// dependent proxies resynced (CRIT: this used to compare names only, so a proxy stayed
// on "header" forever after its provider moved to "query").
func TestProxyAuthHeadersStale(t *testing.T) {
	proxy := func(ingress, ingressIn string, upstream *string) *models.LLMProxy {
		config := models.LLMProxyConfig{}
		if ingress != "" {
			config.Security = &models.SecurityConfig{APIKey: &models.APIKeySecurity{Key: ingress, In: ingressIn}}
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
		ingress, ingressIn        string
		upstream                  string
		wantIngress, wantUpstream bool
	}{
		{
			name:    "both already match",
			proxy:   proxy("x-api-key", "header", header("x-api-key")),
			ingress: "x-api-key", ingressIn: "header", upstream: "x-api-key",
		},
		{
			name:    "ingress name differs",
			proxy:   proxy("API-Key", "header", header("x-api-key")),
			ingress: "x-api-key", ingressIn: "header", upstream: "x-api-key",
			wantIngress: true,
		},
		{
			// Same name, but the provider moved from header to query — the location
			// itself must be recognized as stale, not just the name.
			name:    "ingress location differs",
			proxy:   proxy("x-api-key", "header", header("x-api-key")),
			ingress: "x-api-key", ingressIn: "query", upstream: "x-api-key",
			wantIngress: true,
		},
		{
			name:    "upstream differs",
			proxy:   proxy("x-api-key", "header", header("API-Key")),
			ingress: "x-api-key", ingressIn: "header", upstream: "x-api-key",
			wantUpstream: true,
		},
		{
			// The pre-existing gap: a proxy provisioned before an edit carries both
			// stale names, and fixing only one leaves the other hop broken.
			name:    "both differ",
			proxy:   proxy("API-Key", "header", header("API-Key")),
			ingress: "x-api-key", ingressIn: "header", upstream: "x-api-key",
			wantIngress: true, wantUpstream: true,
		},
		{
			name:    "upstream auth present but unset header",
			proxy:   proxy("x-api-key", "header", nil),
			ingress: "x-api-key", ingressIn: "header", upstream: "x-api-key",
		},
		{
			name:    "proxy without api key security is left alone",
			proxy:   proxy("", "", header("x-api-key")),
			ingress: "x-api-key", ingressIn: "header", upstream: "x-api-key",
		},
		{
			// A provider with no api-key header of its own names nothing to forward,
			// so the upstream hop must not be rewritten to an empty header.
			name:    "provider names no upstream header",
			proxy:   proxy("x-api-key", "header", header("API-Key")),
			ingress: "x-api-key", ingressIn: "header", upstream: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIngress, gotUpstream := proxyAuthHeadersStale(tt.proxy, tt.ingress, tt.ingressIn, tt.upstream)
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

	gotName, gotIn := llmProxyAPIKeySecurity(provisioned)
	require.Equal(t, models.DefaultLLMProxyAPIKeyHeader, gotName)
	require.Equal(t, "header", gotIn)
}

// The sync this gates redeploys every dependent proxy, so it must fire only when a
// provider edit actually changed one of the two headers — never on an edit (name,
// description, policy) that left security untouched, and never on a pre-existing
// mismatch that predates this feature but wasn't part of the current edit.
func TestProviderAuthHeadersChanged(t *testing.T) {
	providerWith := func(key string) *models.LLMProvider {
		return &models.LLMProvider{
			Configuration: models.LLMProviderConfig{
				Security: &models.SecurityConfig{APIKey: &models.APIKeySecurity{Key: key, In: "header"}},
			},
		}
	}

	tests := []struct {
		name              string
		existing, updated *models.LLMProvider
		want              bool
	}{
		{
			name:     "unchanged header",
			existing: providerWith("X-API-Key"),
			updated:  providerWith("X-API-Key"),
			want:     false,
		},
		{
			name:     "renamed header",
			existing: providerWith("API-Key"),
			updated:  providerWith("X-API-Key"),
			want:     true,
		},
		{
			// Both provider configs mirror the default when unset, so this is not a
			// change even though neither names a header explicitly.
			name:     "both fall back to the same default",
			existing: &models.LLMProvider{},
			updated:  &models.LLMProvider{},
			want:     false,
		},
		{
			name:     "nil existing provider",
			existing: nil,
			updated:  providerWith("X-API-Key"),
			want:     true,
		},
		{
			name:     "nil updated provider",
			existing: providerWith("X-API-Key"),
			updated:  nil,
			want:     true,
		},
		{
			// CRIT: a pure location change (same key name, header -> query) used to
			// be invisible here, since the comparison went through a name-only lookup
			// that silently collapsed query locations to the same header default —
			// so switching a provider to query-based auth never triggered a resync.
			name: "same name, location changed from header to query",
			existing: &models.LLMProvider{
				Configuration: models.LLMProviderConfig{
					Security: &models.SecurityConfig{APIKey: &models.APIKeySecurity{Key: "apikey", In: "header"}},
				},
			},
			updated: &models.LLMProvider{
				Configuration: models.LLMProviderConfig{
					Security: &models.SecurityConfig{APIKey: &models.APIKeySecurity{Key: "apikey", In: "query"}},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ProviderAuthHeadersChanged(tt.existing, tt.updated))
		})
	}
}

// SyncDependentProxyAuthHeaders runs detached in its own goroutine (see
// controllers/llm_controller.go), so a second edit to the same provider can commit and
// start its own sync before this one runs. It must resolve headers from the provider's
// current row rather than the snapshot it was handed, or whichever sync finishes last
// can leave proxies on a stale header instead of converging on what's actually stored.
func TestSyncDependentProxyAuthHeaders_RefetchesCurrentProvider(t *testing.T) {
	providerUUID := uuid.New()
	staleSnapshot := &models.LLMProvider{
		UUID: providerUUID,
		Configuration: models.LLMProviderConfig{
			Security: &models.SecurityConfig{APIKey: &models.APIKeySecurity{Key: "B-Header", In: "header"}},
		},
	}
	currentInDB := &models.LLMProvider{
		UUID: providerUUID,
		Configuration: models.LLMProviderConfig{
			Security: &models.SecurityConfig{APIKey: &models.APIKeySecurity{Key: "C-Header", In: "header"}},
		},
	}

	var gotProviderID, gotOuID string
	providerRepo := &repomocks.LLMProviderRepositoryMock{
		GetByUUIDFunc: func(providerID, ouID string) (*models.LLMProvider, error) {
			gotProviderID, gotOuID = providerID, ouID
			return currentInDB, nil
		},
	}
	proxyRepo := &repomocks.LLMProxyRepositoryMock{
		ListByProviderFunc: func(ouID, providerUUID string, limit, offset int) ([]*models.LLMProxy, error) {
			return nil, nil
		},
	}
	svc := &LLMProviderService{providerRepo: providerRepo, proxyRepo: proxyRepo}

	err := svc.SyncDependentProxyAuthHeaders(staleSnapshot, "ou-acme", &LLMProxyService{}, &LLMProxyDeploymentService{})

	require.NoError(t, err)
	require.Equal(t, providerUUID.String(), gotProviderID, "must refetch the provider named in the (possibly stale) snapshot it was handed")
	require.Equal(t, "ou-acme", gotOuID)
}
