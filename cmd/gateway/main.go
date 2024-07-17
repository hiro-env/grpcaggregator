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

func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// デバッグ用に一旦全て許可する
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// For preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	corsMux := allowCORS(mux)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", corsMux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
