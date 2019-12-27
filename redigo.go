package redigo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gozix/viper/v2"
	"github.com/sarulabs/di/v2"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Pool is type alias of redis.Pool
	Pool = redis.Pool
)

// BundleName is default definition name.
const BundleName = "redigo"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: BundleName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cfg *viper.Viper
			if err = ctn.Fill(viper.BundleName, &cfg); err != nil {
				return nil, err
			}

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
					return nil, errors.New(suffix + "host should be set")
				}

				if c.MaxIdle < 0 {
					return nil, errors.New(suffix + "max_idle should be greater or equal to 0")
				}

				if c.MaxActive < 0 {
					return nil, errors.New(suffix + "max_active should be greater or equal to 0")
				}

				if c.IdleTimeout < 0 {
					return nil, errors.New(suffix + "idle_timeout should be greater or equal to 0")
				}

				conf[name] = c
			}

			return NewRegistry(conf), nil
		},
		Close: func(obj interface{}) error {
			return obj.(*Registry).Close()
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{viper.BundleName}
}
