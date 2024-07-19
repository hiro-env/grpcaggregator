package main

import (
	"log"
	"net"
	"os"

	"github.com/hiro-env/grpcaggregator/internal/service"
	"github.com/hiro-env/grpcaggregator/pkg/qiita"
	statsig "github.com/statsig-io/go-sdk"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	agentHost := os.Getenv("DD_AGENT_HOST")
	agentPort := os.Getenv("DD_TRACE_AGENT_PORT")

	tracer.Start(
		tracer.WithAgentAddr(agentHost+":"+agentPort),
		tracer.WithService("grpc-server"),
		tracer.WithEnv("develop"),
	)
	defer tracer.Stop()

	statsigKey := os.Getenv("STATSIG_SERVER_SECRET")
	if statsigKey == "" {
		log.Fatal("statsig server secret is not set!")
	}

	statsig.Initialize(statsigKey)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// インターセプターの導入
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctrace.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpctrace.StreamServerInterceptor()),
	)

	qiitaService := service.NewQiitaService()
	qiita.RegisterQiitaServiceServer(s, qiitaService)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
