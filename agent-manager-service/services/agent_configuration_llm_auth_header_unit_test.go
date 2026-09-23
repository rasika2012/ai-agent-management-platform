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
			// Proxies provisioned before the default was aligned to X-API-Key
			// are still enforced on "API-Key" by their deployed gateway policy,
			// so the stored value has to win over the current default.
			name:  "keeps the legacy header a proxy was provisioned with",
			proxy: proxyWithHeader("API-Key"),
			want:  "API-Key",
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
