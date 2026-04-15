package config

type Config struct {
	App       AppConfig       `toml:"app"`
	Providers ProvidersConfig `toml:"providers"`
	Index     IndexConfig     `toml:"index"`
}

type AppConfig struct {
	ResultLimit int `toml:"result_limit"`
	DebounceMS  int `toml:"debounce_ms"`
}

type ProvidersConfig struct {
	Enabled []string          `toml:"enabled"`
	Web     WebProviderConfig `toml:"web"`
}

type WebProviderConfig struct {
	TimeoutMS int `toml:"timeout_ms"`
}

type IndexConfig struct {
	Path  string   `toml:"path"`
	Watch bool     `toml:"watch"`
	Roots []string `toml:"roots"`
}
