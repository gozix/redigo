package redigo

import (
	"errors"
	"net"
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

	// config is redis configuration struct.
	config struct {
		Host        string
		Port        string
		Password    string
		MaxIdle     int
		MaxActive   int
		IdleTimeout time.Duration
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
			var v *viper.Viper
			if err = ctn.Fill(viper.BundleName, &v); err != nil {
				return nil, err
			}

			// set default
			v.SetDefault("redis.port", "6379")
			v.SetDefault("redis.max_idle", 3)
			v.SetDefault("redis.idle_timeout", 240*time.Second)

			var cfg = config{
				Host:        v.GetString("redis.host"),
				Port:        v.GetString("redis.port"),
				Password:    v.GetString("redis.password"),
				MaxIdle:     v.GetInt("redis.max_idle"),
				MaxActive:   v.GetInt("redis.max_active"),
				IdleTimeout: v.GetDuration("redis.idle_timeout"),
			}

			// validating
			if cfg.Host == "" {
				return nil, errors.New("redis.host should be set")
			}

			if cfg.MaxIdle < 0 {
				return nil, errors.New("redis.max_idle should be greater then 0")
			}

			if cfg.MaxActive < 0 {
				return nil, errors.New("redis.max_active should be greater or equal to 0")
			}

			if cfg.IdleTimeout < 0 {
				return nil, errors.New("redis.idle_timeout should be greater or equal to 0")
			}

			var options []redis.DialOption
			if cfg.Password != "" {
				options = append(options, redis.DialPassword(cfg.Password))
			}

			var pool = &redis.Pool{
				MaxIdle:     cfg.MaxIdle,
				IdleTimeout: cfg.IdleTimeout,
				Dial: func() (redis.Conn, error) {
					return redis.Dial(
						"tcp",
						net.JoinHostPort(
							cfg.Host,
							cfg.Port,
						),
						options...,
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
