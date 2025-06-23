package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
)

func setupRedis7(config *Config) (*redis.Client, error) {
	var (
		ctx context.Context
		rdb *redis.Client
	)
	ctx = context.Background()

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	addr := fmt.Sprintf("%s:%v", config.Address, config.Port)
	err := backoff.Retry(func() error {
		rdb = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     config.Password,
			DB:           config.DB,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolSize:     config.PoolSize,
		})
		err := rdb.Ping(ctx).Err()
		if err != nil {
			return err
		}
		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Ping to redis %s success", addr)

	return rdb, nil
}

func setupRedis7Cluster(config *ClusterConfig) (*redis.ClusterClient, error) {
	var (
		ctx context.Context
		rdb *redis.ClusterClient
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	ctx = context.Background()

	err := backoff.Retry(func() error {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			NewClient: func(opt *redis.Options) *redis.Client {
				return redis.NewClient(opt)
			},
			Addrs:        config.Address,
			MaxRedirects: config.MaxRetries,
			Password:     config.Password,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolSize:     config.PoolSize,
		})

		err := rdb.Ping(ctx).Err()
		if err != nil {
			log.Error().Msgf("redis cluster ping failed: %v", err.Error())
			return err
		}

		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	log.Info().
		Strs("addrs", config.Address).
		Msg("Ping to redis cluster success")

	return rdb, nil
}

func setupRedis7Sentinel(config *SentinelConfig) (*redis.Client, error) {
	var (
		ctx context.Context
		rdb *redis.Client
	)
	log.Info().Msgf("address %v master %v pass %v", config.Address, config.MasterName, config.Password)
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	ctx = context.Background()

	err := backoff.Retry(func() error {
		rdb = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       config.MasterName,
			SentinelAddrs:    config.Address,
			SentinelPassword: config.Password,
			Password:         config.Password,
			MaxRetries:       config.MaxRetries,
			DialTimeout:      config.DialTimeout,
			ReadTimeout:      config.ReadTimeout,
			WriteTimeout:     config.WriteTimeout,
			PoolSize:         config.PoolSize,
		})

		err := rdb.Ping(ctx).Err()
		if err != nil {
			log.Error().Msgf("err %s", err.Error())
			return err
		}

		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	log.Info().
		Strs("addrs", config.Address).
		Msg("Ping to redis sentinel success")

	return rdb, nil
}
