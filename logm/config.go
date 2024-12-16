package logm

import (
	"encoding"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/sirupsen/logrus"

	"github.com/ibrt/golang-modules/cfgm"
)

var (
	_ encoding.TextUnmarshaler = (*LogConfigLogrusOutput)(nil)
	_ cfgm.Config              = (*LogConfig)(nil)
	_ LogConfigMixin           = (*LogConfig)(nil)
)

// LogConfigLogrusOutput describes the acceptable values for [LogConfig.LogrusOutput].
type LogConfigLogrusOutput string

// Known LogConfigLogrusOutput values.
const (
	LogConfigLogrusOutputDisabled LogConfigLogrusOutput = cfgm.DisabledValue
	LogConfigLogrusOutputHuman    LogConfigLogrusOutput = "human"
	LogConfigLogrusOutputJSON     LogConfigLogrusOutput = "json"
)

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (l *LogConfigLogrusOutput) UnmarshalText(text []byte) error {
	switch v := LogConfigLogrusOutput(text); v {
	case LogConfigLogrusOutputDisabled, LogConfigLogrusOutputHuman, LogConfigLogrusOutputJSON:
		*l = v
		return nil
	default:
		return errorz.Errorf("invalid value for LogConfigLogrusOutput: '%s'", v)
	}
}

// LogConfigMixin describes the module configuration.
type LogConfigMixin interface {
	cfgm.Config
	GetLogConfig() *LogConfig
}

// LogConfig describes the module configuration.
type LogConfig struct {
	HoneycombAPIKey     string                `env:"LOG_HONEYCOMB_API_KEY,required"`
	HoneycombDataset    string                `env:"LOG_HONEYCOMB_DATASET,required"`
	HoneycombSampleRate uint                  `env:"LOG_HONEYCOMB_SAMPLE_RATE,required" validate:"required,min=1"`
	LogrusOutput        LogConfigLogrusOutput `env:"LOG_LOGRUS_OUTPUT,required"`
	LogrusLevel         logrus.Level          `env:"LOG_LOGRUS_LEVEL,required"`
}

// Config implements the [cfgm.Config] interface.
func (c *LogConfig) Config() {
	// intentionally empty
}

// GetLogConfig implements the [LogConfigMixin] interface.
func (c *LogConfig) GetLogConfig() *LogConfig {
	return c
}
