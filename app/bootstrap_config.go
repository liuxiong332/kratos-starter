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

	flag.Parse()

	copyIfNotEmpty(configPath, &config.ConfigPath)
	copyIfNotEmpty(address, &config.ConsulAddress)
	copyIfNotEmpty(token, &config.ConsulToken)
	copyIfNotEmpty(vaultToken, &config.VaultToken)
	copyIfNotEmpty(vaultPrefix, &config.VaultPrefix)
}
