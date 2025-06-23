package redis

import (
	redis7 "github.com/redis/go-redis/v9"
)

// NewRedis7Client new go-redis client
func NewRedis7Client(config *Config) (*redis7.Client, error) {
	return setupRedis7(config)
}

// NewRedis7Cluster new go-redis cluster
func NewRedis7Cluster(config *ClusterConfig) (*redis7.ClusterClient, error) {
	return setupRedis7Cluster(config)
}

// NewRedis67entinel new go-redis Sentinel
func NewRedis67entinel(config *SentinelConfig) (*redis7.Client, error) {
	return setupRedis7Sentinel(config)
}
