/**
* @description:
* @author
* @date 2026-03-25 23:27:11
* @version 1.0
*
* Change Logs:
* Date           Author       Notes
*
 */

package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"pxsemic.com/simplebank/mail"
	"pxsemic.com/simplebank/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"pxsemic.com/simplebank/api"
	db "pxsemic.com/simplebank/db/sqlc"
	"pxsemic.com/simplebank/gapi"
	"pxsemic.com/simplebank/pb"
	"pxsemic.com/simplebank/util"

	_ "github.com/lib/pq"
	_ "pxsemic.com/simplebank/doc/statik"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("can't load config app.env. err:")
	}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02 15:04:05.000", // 设置时间格式
		})
		log.Info().Msgf("start app in %s mode", config.Environment)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open db")
	}
	store := db.NewStore(conn)
	runMigrationDb(config.MigrationURL, config.DBSource)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)

	//runHttpServer(config, store)
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewYMailSend(config.EmailSenderName, config.EmailSenderEmail, config.EmailSenderPassword)

	taskProcessor := worker.NewRedisTaskProcessor(store, redisOpt, mailer)
	log.Info().Msg("start task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

// runMigrationDb 运行数据库迁移
func runMigrationDb(migrationURL string, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migrate instance")
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}
	log.Info().Msg("db migrated successfully")
}

// runGrpcServer 运行gRPC服务器
func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(store, config, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create gapi server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen gRPC server")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start gRPC server")
	}
}

// runGatewayServer 运行gRPC网关服务器
func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(store, config, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create gapi server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	handler := gapi.HttpLogger(mux)
	httpServer := &http.Server{
		Handler: handler,
		Addr:    config.HTTPServerAddress,
	}

	log.Printf("start http gateway server at %s", httpServer.Addr)
	if err = httpServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start gRPC gateway server")
	}
}
func runHttpServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
