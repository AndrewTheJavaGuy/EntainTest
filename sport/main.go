package main

import (
	"database/sql"
	"flag"
	"github.com/AndrewTheJavaGuy/entain/sport/db"
	"github.com/AndrewTheJavaGuy/entain/sport/proto/sport"
	"github.com/AndrewTheJavaGuy/entain/sport/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var (
	grpcEndpoint = flag.String("Sports-endpoint", "localhost:9001", "Sports server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running Sports server: %s", err)
	}
}

func run() error {

	conn, err := net.Listen("tcp", ":9001")
	if err != nil {
		return err
	}

	sportsDB, err := sql.Open("sqlite3", "./db/sports.db")
	if err != nil {
		return err
	}

	sportsRepo := db.NewSportsRepo(sportsDB)
	if err := sportsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sport.RegisterSportServer(
		grpcServer,
		service.NewSportsService(
			sportsRepo,
		),
	)

	log.Infof("gRPC server listening on: %s", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
