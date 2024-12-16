package logm

import (
	"encoding"
	"fmt"

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

// String implements the [fmt.Stringer] interface.
func (l *LogConfigLogrusOutput) String() string {
	return string(*l)
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

// ToEnv converts the config to an env map.
func (c *LogConfig) ToEnv(prefix string) map[string]string {
	return map[string]string{
		prefix + "LOG_HONEYCOMB_API_KEY":     c.HoneycombAPIKey,
		prefix + "LOG_HONEYCOMB_DATASET":     c.HoneycombDataset,
		prefix + "LOG_HONEYCOMB_SAMPLE_RATE": fmt.Sprintf("%v", c.HoneycombSampleRate),
		prefix + "LOG_LOGRUS_OUTPUT":         c.LogrusOutput.String(),
		prefix + "LOG_LOGRUS_LEVEL":          c.LogrusLevel.String(),
	}
}

// Config implements the [cfgm.Config] interface.
func (c *LogConfig) Config() {
	// intentionally empty
}

// GetLogConfig implements the [LogConfigMixin] interface.
func (c *LogConfig) GetLogConfig() *LogConfig {
	return c
}
