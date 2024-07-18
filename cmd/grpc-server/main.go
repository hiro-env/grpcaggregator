package main

import (
	"log"
	"net"
	"os"

	"github.com/hiro-env/grpcaggregator/internal/service"
	"github.com/hiro-env/grpcaggregator/pkg/qiita"
	statsig "github.com/statsig-io/go-sdk"
	"google.golang.org/grpc"
)

func main() {
	statsigKey := os.Getenv("STATSIG_SERVER_SECRET")
	if statsigKey == "" {
		log.Fatal("statsig server secret is not set!")
	}

	statsig.Initialize(statsigKey)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	qiitaService := service.NewQiitaService()
	qiita.RegisterQiitaServiceServer(s, qiitaService)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
