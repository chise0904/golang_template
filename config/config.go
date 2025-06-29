package config

import (
	db "github.com/chise0904/golang_template/pkg/database/gorm"
	// "github.com/chise0904/golang_template/pkg/database/redis"
	"github.com/chise0904/golang_template/pkg/grpc"

	// jetstream_client "github.com/chise0904/golang_template/pkg/messaging/nats/jetstream"

	// "github.com/chise0904/golang_template/pkg/recommender/gorse"
	web "github.com/chise0904/golang_template/pkg/web"
	"github.com/chise0904/golang_template/pkg/zlog"
	"github.com/chise0904/golang_template/service"

	"go.uber.org/fx"
)

type Config struct {
	fx.Out
	Log       *zlog.Config `mapstructure:"log"`
	Database  *db.Config   `mapstructure:"database"`
	Grpc      *grpc.Config `mapstructure:"grpc"`
	WebConfig *web.Config  `mapstructure:"web"`
	// GorseConfig     *gorse.Config            `mapstructure:"gorse"`
	// RedisCluster *redis.ClusterConfig `mapstructure:"redis"`
	// JetstreamConfig *jetstream_client.Config `mapstructure:"jetstream"`

	// IdentityWebServiceConfig *service.WebServiceConfig `mapstructure:"webService"` // connect to identity web static resource
	AccessConfig *service.AccessTokenConfig `mapstructure:"access"`
	// OAuthConfig              *service.OAuthConfig         `mapstructure:"oauth"`
	// NotifyServiceConfig      *service.NotifyServiceConfig `mapstructure:"notify_service"`
}
