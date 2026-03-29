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
	"database/sql"
	"log"
	"net"
	"net/http"

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
		log.Fatal("can't load config app.env. err:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot open db", err)
	}
	store := db.NewStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)

	//runHttpServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("can't create gapi server", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("failed to listen http server", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("failed to start gRPC server", err)
	}
}
func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("can't create gapi server", err)
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("failed to register gRPC gateway handler", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("cannot create statik fs: %v", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	httpMux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("failed to listen gRPC gateway server", err)
	}
	log.Printf("start gRPC gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, httpMux)
	if err != nil {
		log.Fatal("failed to start gRPC gateway server", err)
	}
}
func runHttpServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
