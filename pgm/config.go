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
	PostgresURL string `env:"POSTGRES_URL,required" validate:"required,url"`
}

// Config implements the [cfgm.Config] interface.
func (c *PGConfig) Config() {
	// intentionally empty
}

// GetPGConfig implements the [PGConfigMixin] interface.
func (c *PGConfig) GetPGConfig() *PGConfig {
	return c
}
