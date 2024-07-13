package main

import (
	"log"
	"net"

	"github.com/hiro-env/grpcaggregator/internal/service"
	"github.com/hiro-env/grpcaggregator/pkg/qiita"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	qiita.RegisterQiitaServiceServer(s, &service.QiitaService{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
