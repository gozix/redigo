// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package redigo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gozix/di"
	"github.com/gozix/glue/v3"
	"github.com/gozix/viper/v2"
)

// Bundle implements the glue.Bundle interface.
type Bundle struct{}

// BundleName is default definition name.
const BundleName = "redigo"

// Bundle implements the glue.Bundle interface.
var _ glue.Bundle = (*Bundle)(nil)

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder di.Builder) error {
	return builder.Provide(b.provideRegistry)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{
		viper.BundleName,
	}
}

func (b *Bundle) provideRegistry(cfg *viper.Viper) (*Registry, func() error, error) {
	// use this is hack, not UnmarshalKey
	// see https://github.com/spf13/viper/issues/188
	var (
		keys = cfg.Sub(configKey).AllKeys()
		conf = make(Configs, len(keys))
	)

	for _, key := range keys {
		var name = strings.Split(key, ".")[0]
		if _, ok := conf[name]; ok {
			continue
		}

		var suffix = fmt.Sprintf("%s.%s.", configKey, name)

		cfg.SetDefault(suffix+"port", "6379")
		cfg.SetDefault(suffix+"max_idle", 3)
		cfg.SetDefault(suffix+"idle_timeout", 240*time.Second)

		var c = Config{
			Host:        cfg.GetString(suffix + "host"),
			Port:        cfg.GetString(suffix + "port"),
			DB:          cfg.GetInt(suffix + "db"),
			Password:    cfg.GetString(suffix + "password"),
			MaxIdle:     cfg.GetInt(suffix + "max_idle"),
			MaxActive:   cfg.GetInt(suffix + "max_active"),
			IdleTimeout: cfg.GetDuration(suffix + "idle_timeout"),
		}

		// validating
		if c.Host == "" {
			return nil, nil, errors.New(suffix + "host should be set")
		}

		if c.DB < 0 {
			return nil, nil, errors.New(suffix + "db should be greater or equal to 0")
		}

		if c.MaxIdle < 0 {
			return nil, nil, errors.New(suffix + "max_idle should be greater or equal to 0")
		}

		if c.MaxActive < 0 {
			return nil, nil, errors.New(suffix + "max_active should be greater or equal to 0")
		}

		if c.IdleTimeout < 0 {
			return nil, nil, errors.New(suffix + "idle_timeout should be greater or equal to 0")
		}

		conf[name] = c
	}

	var (
		registry = NewRegistry(conf)
		closer   = func() error {
			return registry.Close()
		}
	)

	return registry, closer, nil
}
