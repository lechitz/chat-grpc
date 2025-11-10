// Package observability exposes helper functions shared across tracing and metrics setup.
package observability

import (
	"net/url"
	"strings"
)

// ParseHeaders breaks a comma-separated list of k=v pairs into a string map suitable for OTEL options.
func ParseHeaders(headersStr string) map[string]string {
	result := make(map[string]string)
	if strings.TrimSpace(headersStr) == "" {
		return result
	}

	pairs := strings.Split(headersStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		if key == "" {
			continue
		}
		result[key] = value
	}

	return result
}

// NormalizeEndpoint ensures an OTLP endpoint string is a valid URL, accepting both host:port and full URLs.
func NormalizeEndpoint(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", nil
	}

	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// ExportHostFromEndpoint returns the host:port to be used by OTLP gRPC exporters.
func ExportHostFromEndpoint(endpoint string) (string, error) {
	normalized, err := NormalizeEndpoint(endpoint)
	if err != nil {
		// return raw endpoint as fallback
		return endpoint, err
	}

	if strings.HasPrefix(normalized, "http://") || strings.HasPrefix(normalized, "https://") {
		if u, err := url.Parse(normalized); err == nil && u.Host != "" {
			return u.Host, nil
		}
	}

	return normalized, nil
}
