package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hiro-env/grpcaggregator/pkg/qiita"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := qiita.RegisterQiitaServiceHandlerFromEndpoint(ctx, mux, "grpc-server:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
