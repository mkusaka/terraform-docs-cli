package registry

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewClient_UsesProxyFromEnvironment(t *testing.T) {
	proxyURL := "http://127.0.0.1:18080"
	oldHTTPProxy := os.Getenv("HTTP_PROXY")
	oldHTTPSProxy := os.Getenv("HTTPS_PROXY")
	defer func() {
		_ = os.Setenv("HTTP_PROXY", oldHTTPProxy)
		_ = os.Setenv("HTTPS_PROXY", oldHTTPSProxy)
	}()
	_ = os.Setenv("HTTP_PROXY", proxyURL)
	_ = os.Setenv("HTTPS_PROXY", proxyURL)

	c, err := NewClient(Config{BaseURL: "https://registry.terraform.io", Timeout: 5 * time.Second}, nil)
	if err != nil {
		t.Fatal(err)
	}

	transport, ok := c.httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected transport type: %T", c.httpClient.Transport)
	}
	if transport.Proxy == nil {
		t.Fatalf("expected proxy function to be set")
	}

	req, err := http.NewRequest(http.MethodGet, "https://registry.terraform.io/v2/providers/hashicorp/aws", nil)
	if err != nil {
		t.Fatal(err)
	}
	gotProxy, err := transport.Proxy(req)
	if err != nil {
		t.Fatal(err)
	}
	if gotProxy == nil || gotProxy.String() != proxyURL {
		t.Fatalf("expected proxy %s, got %v", proxyURL, gotProxy)
	}
}

func TestNewClient_InvalidBaseURLWithoutSchemeOrHostReturnsConfigError(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantMsg string
	}{
		{name: "missing scheme", baseURL: "registry.terraform.io", wantMsg: "scheme and host are required"},
		{name: "missing host", baseURL: "https:///v2", wantMsg: "scheme and host are required"},
		{name: "unsupported scheme", baseURL: "ftp://registry.terraform.io", wantMsg: "scheme must be http or https"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(Config{BaseURL: tt.baseURL, Timeout: 5 * time.Second}, nil)
			if err == nil {
				t.Fatalf("expected error for base url %q", tt.baseURL)
			}

			var cfgErr *ConfigError
			if !errors.As(err, &cfgErr) {
				t.Fatalf("expected ConfigError, got %T (%v)", err, err)
			}
			if !strings.Contains(cfgErr.Error(), tt.wantMsg) {
				t.Fatalf("unexpected error message: %s", cfgErr.Error())
			}
		})
	}
}

func TestResolve_PreservesBasePathPrefixForAbsoluteAPIPaths(t *testing.T) {
	c, err := NewClient(Config{BaseURL: "https://example.com/registry", Timeout: 5 * time.Second}, nil)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.resolve("/v2/providers/hashicorp/aws?include=provider-versions")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://example.com/registry/v2/providers/hashicorp/aws?include=provider-versions"
	if got != want {
		t.Fatalf("unexpected resolved URL\nwant: %s\ngot:  %s", want, got)
	}
}

func TestResolve_RootBasePathStillUsesRoot(t *testing.T) {
	c, err := NewClient(Config{BaseURL: "https://registry.terraform.io", Timeout: 5 * time.Second}, nil)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.resolve("/v2/providers/hashicorp/aws")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://registry.terraform.io/v2/providers/hashicorp/aws"
	if got != want {
		t.Fatalf("unexpected resolved URL\nwant: %s\ngot:  %s", want, got)
	}
}
