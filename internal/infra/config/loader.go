package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

func Load(path string) (*Config, error) {
	cfg := Default()

	if path == "" {
		return cfg, nil
	}

	if _, err := os.Stat(path); err != nil {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func Default() *Config {
	return &Config{
		App: AppConfig{
			ResultLimit: 20,
			DebounceMS:  250,
		},
		Providers: ProvidersConfig{
			Enabled: []string{"local"},
			Web: WebProviderConfig{
				TimeoutMS: 1200,
			},
		},
		Index: IndexConfig{
			Path:  "./data/index",
			Watch: false,
			Roots: []string{},
		},
	}
}
