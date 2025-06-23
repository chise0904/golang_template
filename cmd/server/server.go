package server

import (
	"context"
	"os"
	"time"

	"github.com/chise0904/golang_template/pkg/config"
	db "github.com/chise0904/golang_template/pkg/database/gorm"

	// "github.com/chise0904/golang_template/pkg/database/redis"
	// jetstream_client "github.com/chise0904/golang_template/pkg/messaging/nats/jetstream"

	// "github.com/chise0904/golang_template/pkg/recommender/gorse"
	web "github.com/chise0904/golang_template/pkg/web/echo"

	configs "github.com/chise0904/golang_template/config"
	"github.com/chise0904/golang_template/delivery"

	// grpc_delivery "github.com/chise0904/golang_template/delivery/grpc"
	"github.com/spf13/cobra"

	repo_impl "github.com/chise0904/golang_template/repository/impl"
	// identity_service "github.com/chise0904/golang_template/service/identity_service"
	// jwt_auth "github.com/chise0904/golang_template/service/jwt_auth_service"
	// notify "github.com/chise0904/golang_template/service/notify_service"

	"github.com/chise0904/golang_template/pkg/zlog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func ServerCmd() *cobra.Command {

	return &cobra.Command{
		Use: "server",
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	cfg := configs.Config{}
	err := config.LoadConfig(os.Getenv("CONFIG_PATH"), &cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	zlog.Setup(cfg.Log)
	app := fx.New(
		fx.WithLogger(zlog.FxLogger()),
		fx.Supply(cfg), // 	已經建構好的變數，如 config, logger
		fx.Provide( // fx 會自動注入參數，並建立實例
			db.NewConnection,
			// jetstream_client.NewJetStream,
			// redis.NewRedis7Cluster,
			web.NewEcho,
			repo_impl.NewIdentityRepo,
			// jwt_auth.NewJWTAuthService,
			// notify.NewNotifyService,
			// identity_service.NewIdentityService,
			// gorse.NewGorseClient,
		),
		fx.Invoke( // 用於啟動服務、註冊 handler 等
			delivery.SetIdentityDelivery,
			// grpc_delivery.RunStorageServiceGrpcHandler,
			// notify.RunUserNotificationConsumer,
		),
	)

	log.Info().Msg("launch identity-server")
	app.Run()

	log.Info().Msg("main: shutting down identity-server...")
	exitCode := 0
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Error().Msgf("main: server shutdown error: %v", err)
		exitCode++
	} else {
		log.Info().Msg("main: gracefully stopped")
	}
	os.Exit(exitCode)
}
