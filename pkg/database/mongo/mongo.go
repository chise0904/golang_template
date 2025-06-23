package mongo

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// SetupMongo .
func SetupMongo(cfg *Config) (*mongo.Database, error) {
	ctx := context.Background()
	var uri string
	credential := options.Credential{Username: cfg.Username, Password: cfg.Password}
	if cfg.Port == "" {
		credential = options.Credential{AuthMechanism: "SCRAM-SHA-1", Username: cfg.Username, Password: cfg.Password, AuthSource: "admin"}
		uri = fmt.Sprintf("mongodb+srv://%s", cfg.Host)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
	}

	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Debug().Msgf("mongo cmd: %+v", evt.Command)
		},
	}

	opts := options.Client().
		ApplyURI(uri).
		SetAuth(credential)
	if cfg.Monitor == "debug" {
		opts = opts.SetMonitor(cmdMonitor)
	}
	retry := true
	opts.RetryWrites = &retry
	opts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	// 連接 mongo
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// 檢查連接
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(cfg.Database), nil
}
