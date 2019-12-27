package redigo

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// DEFAULT is default connection name.
const DEFAULT = "default"

// ConfigKey is root config key.
const configKey = "redis"

type (
	// Config is registry configuration item.
	Config struct {
		Host        string        `json:"host"`
		Port        string        `json:"port"`
		DB          int           `json:"db"`
		Password    string        `json:"password"`
		MaxIdle     int           `json:"max_idle"`
		MaxActive   int           `json:"max_active"`
		IdleTimeout time.Duration `json:"idle_timeout"`
	}

	// Configs is registry configurations.
	Configs map[string]Config

	// Registry is database connection registry.
	Registry struct {
		mux   sync.Mutex
		pools map[string]*redis.Pool
		conf  Configs
	}
)

var (
	// ErrUnknownConnection is error triggered when connection with provided name not founded.
	ErrUnknownConnection = errors.New("unknown connection")
)

// NewRegistry is registry constructor.
func NewRegistry(conf Configs) *Registry {
	return &Registry{
		pools: make(map[string]*redis.Pool, 1),
		conf:  conf,
	}
}

// Close is method for close connections.
func (r *Registry) Close() (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for key, pool := range r.pools {
		if err = pool.Close(); err != nil {
			return err
		}

		delete(r.pools, key)
	}

	return nil
}

// Connection is default connection getter.
func (r *Registry) Connection() (*redis.Pool, error) {
	return r.ConnectionWithName(DEFAULT)
}

// ConnectionWithName is connection getter by name.
func (r *Registry) ConnectionWithName(name string) (_ *redis.Pool, err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var pool, initialized = r.pools[name]
	if initialized {
		return pool, nil
	}

	var cfg, exists = r.conf[name]
	if !exists {
		return nil, ErrUnknownConnection
	}

	var options = []redis.DialOption{
		redis.DialDatabase(cfg.DB),
	}
	if cfg.Password != "" {
		options = append(options, redis.DialPassword(cfg.Password))
	}

	pool = &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
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

	r.pools[name] = pool

	return pool, nil
}
