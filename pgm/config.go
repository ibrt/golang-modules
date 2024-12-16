package pgm

import (
	"github.com/ibrt/golang-modules/cfgm"
)

// PGConfigMixin describes the module configuration.
type PGConfigMixin interface {
	cfgm.Config
	GetPGConfig() *PGConfig
}

// PGConfig describes the module configuration.
type PGConfig struct {
	PostgresURL string `env:"PG_POSTGRES_URL,required" validate:"required,url"`
}

// ToEnv converts the config to an env map.
func (c *PGConfig) ToEnv(prefix string) map[string]string {
	return map[string]string{
		prefix + "PG_POSTGRES_URL": c.PostgresURL,
	}
}

// Config implements the [cfgm.Config] interface.
func (c *PGConfig) Config() {
	// intentionally empty
}

// GetPGConfig implements the [PGConfigMixin] interface.
func (c *PGConfig) GetPGConfig() *PGConfig {
	return c
}
