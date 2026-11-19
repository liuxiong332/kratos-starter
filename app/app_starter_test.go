package app

import "testing"

func TestVaultConfigPath(t *testing.T) {
	cases := []struct {
		prefix  string
		appName string
		want    string
	}{
		{"", "my-worker", "secret/my-worker"},
		{"   ", "my-worker", "secret/my-worker"},
		{"secret", "my-worker", "secret/my-worker"},
		{"vega-secret", "my-worker", "vega-secret/my-worker"},
		{"/vega-secret/", "my-worker", "vega-secret/my-worker"},
		{" vega-secret/data ", "my-worker", "vega-secret/data/my-worker"},
		{"vega-secret", "/my-worker", "vega-secret/my-worker"},
	}

	for _, c := range cases {
		if got := vaultConfigPath(c.prefix, c.appName); got != c.want {
			t.Errorf("vaultConfigPath(%q, %q) = %q, want %q", c.prefix, c.appName, got, c.want)
		}
	}
}

func TestResolveServiceName(t *testing.T) {
	cases := []struct {
		serviceName string
		appName     string
		want        string
	}{
		{"", "my-worker", "my-worker"},
		{"  ", "my-worker", "my-worker"},
		{"my-worker", "my-worker", "my-worker"},
		{"my-worker-canary", "my-worker", "my-worker-canary"},
		{" my-worker-canary ", "my-worker", "my-worker-canary"},
	}

	for _, c := range cases {
		if got := resolveServiceName(c.serviceName, c.appName); got != c.want {
			t.Errorf("resolveServiceName(%q, %q) = %q, want %q", c.serviceName, c.appName, got, c.want)
		}
	}
}
