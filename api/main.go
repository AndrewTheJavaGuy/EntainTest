package main

import (
	"context"
	"flag"
	"git.neds.sh/matty/entain/api/proto/sport"
	"net/http"

	"git.neds.sh/matty/entain/api/proto/racing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	apiEndpoint  = flag.String("api-endpoint", "localhost:8000", "API endpoint")
	racingEndpoint = flag.String("racing-endpoint", "localhost:9000", "Racing server endpoint")
	sportsEndpoint = flag.String("sport-endpoint", "localhost:9001", "Sports server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running api server: %s", err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	if err := racing.RegisterRacingHandlerFromEndpoint(
		ctx,
		mux,
		*racingEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	if err := sport.RegisterSportHandlerFromEndpoint(
		ctx,
		mux,
		*sportsEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}
	log.Infof("API server listening on: %s", *apiEndpoint)

	return http.ListenAndServe(*apiEndpoint, mux)
}
