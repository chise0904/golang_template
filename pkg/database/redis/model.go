package redis

import (
	"time"
)

// Config redis config
type Config struct {
	Address      string        `mapstructure:"addr"`
	Port         int           `mapstructure:"port"`
	DB           int           `mapstructure:"db"`
	Password     string        `mapstructure:"password"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// ClusterConfig redis cluster config
type ClusterConfig struct {
	Address      []string      `mapstructure:"addrs"`
	Password     string        `mapstructure:"password"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// SentinelConfig redis sentinel config
type SentinelConfig struct {
	Address      []string      `mapstructure:"addrs"`
	Password     string        `mapstructure:"password"`
	MasterName   string        `mapstructure:"master_name"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}
