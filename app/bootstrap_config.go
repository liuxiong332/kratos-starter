package app

import (
	"flag"
	"os"
)

type BootstrapConfig struct {
	ConfigPath    string
	ConsulAddress string
	ConsulToken   string
	ConsulTags    string
	VaultToken    string
	// VaultPrefix is the vault mount/path prefix the app config is read from,
	// e.g. "secret" or "vega-secret". Empty means the default "secret".
	VaultPrefix string
	// ServiceName is the consul service name the app registers itself under.
	// Empty means the app name given to NewApp, which is also the name the
	// config is read by (config/{appName} and {vaultPrefix}/{appName}).
	ServiceName string
}

func copyIfNotEmpty(str *string, target *string) {
	if *str != "" {
		*target = *str
	}
}

func ParseBootstrapConfigEnv() *BootstrapConfig {
	config := BootstrapConfig{
		os.Getenv("APP_CONFIG_PATH"),
		os.Getenv("APP_CONSUL_ADDRESS"),
		os.Getenv("APP_CONSUL_TOKEN"),
		os.Getenv("APP_CONSUL_TAGS"),
		os.Getenv("APP_VAULT_TOKEN"),
		os.Getenv("APP_VAULT_PREFIX"),
		os.Getenv("APP_SERVICE_NAME"),
	}

	if !flag.Parsed() {
		ParseBootstrapConfigFlag(&config)
	}

	return &config
}

func ParseBootstrapConfigFlag(config *BootstrapConfig) {
	configPath := flag.String("config_path", "", "Config path")
	address := flag.String("consul_address", "", "Consul Address like localhost:8500")
	token := flag.String("consul_token", "", "Consul Token")
	vaultToken := flag.String("vault_token", "", "Vault Token")
	vaultPrefix := flag.String("vault_prefix", "", "Vault config path prefix, defaults to secret")
	serviceName := flag.String("service_name", "", "Consul service name to register, defaults to the app name")

	flag.Parse()

	copyIfNotEmpty(configPath, &config.ConfigPath)
	copyIfNotEmpty(address, &config.ConsulAddress)
	copyIfNotEmpty(token, &config.ConsulToken)
	copyIfNotEmpty(vaultToken, &config.VaultToken)
	copyIfNotEmpty(vaultPrefix, &config.VaultPrefix)
	copyIfNotEmpty(serviceName, &config.ServiceName)
}
