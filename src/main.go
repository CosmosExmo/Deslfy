package main

import (
	"context"
	"database/sql"
	"desly/api"
	db "desly/db/sqlc"
	_ "desly/doc/statik"
	"desly/gapi"
	"desly/mail"
	"desly/pb"
	"desly/util"
	"desly/worker"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == util.EnvironmentDevelopment {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	//Run db migration
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	host := config.RedisServerAddress
	port := ""
	redisAddr := host + ":" + port;

	//Get the k8n cluster redis address and port
	if(config.Environment == util.EnvironmentProduction) {
		if val := os.Getenv("REDIS_HOST"); val != "" {
			host = val
		}
		if val := string(os.Getenv("REDIS_PORT")); val != "" {
				port = val
		}
		redisAddr = host + ":" + port;
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	
	runTaskProcessor(config, redisOpt, store)
	go runGrpcServer(config, store, taskDistributor)
	go runGatewayServer(config, store, taskDistributor)
	runGinServer(config, store)
}

func runTaskProcessor(config util.Config,redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewEmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msgf("start task processor at: %s", redisOpt.Addr)
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	log.Info().Msgf("start HTTP GIN server")
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterDeslfyServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterDeslfyHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.NewWithNamespace("api_docs")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load api_docs statik files")
	}

	swaggerHandler := http.StripPrefix("/docs/", http.FileServer(statikFS))
	mux.Handle("/docs/", swaggerHandler)

	listener, err := net.Listen("tcp", config.GRPCGatewayServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}
