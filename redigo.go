package redigo

import (
	"errors"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gozix/viper"
	"github.com/sarulabs/di"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Pool is type alias of redis.Pool
	Pool = redis.Pool

	// redisConf is logger configuration struct.
	redisConf struct {
		Host        string
		Port        string
		MaxIdle     int           `mapstructure:"max_idle"`
		MaxActive   int           `mapstructure:"max_active"`
		IdleTimeout time.Duration `mapstructure:"idle_timeout"`
	}
)

// BundleName is default definition name.
const BundleName = "redigo"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Key implements the glue.Bundle interface.
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

			var conf redisConf
			if err = cfg.UnmarshalKey("redis", &conf); err != nil {
				return nil, err
			}

			if conf.MaxIdle < 0 {
				return nil, errors.New("redis.MaxIdle should be greater then 0")
			}

			// set default
			if conf.MaxIdle == 0 {
				conf.MaxIdle = 3
			}

			if conf.MaxActive < 0 {
				return nil, errors.New("redis.MaxActive should be greater or equal to 0")
			}

			if conf.IdleTimeout < 0 {
				return nil, errors.New("redis.IdleTimeout should be greater or equal to 0")
			}

			if conf.IdleTimeout == 0 {
				conf.IdleTimeout = 240 * time.Second
			}

			var pool = &redis.Pool{
				MaxIdle:     conf.MaxIdle,
				IdleTimeout: conf.IdleTimeout,
				Dial: func() (redis.Conn, error) {
					return redis.Dial(
						"tcp",
						net.JoinHostPort(
							cfg.GetString("redis.host"),
							cfg.GetString("redis.port"),
						),
					)
				},
			}

			var conn = pool.Get()
			defer conn.Close()

			if _, err = conn.Do("PING"); err != nil {
				return nil, err
			}

			return pool, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*redis.Pool).Close()
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{viper.BundleName}
}
