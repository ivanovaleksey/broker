package main

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"github.com/ivanovaleksey/broker/src/broker"
	"github.com/ivanovaleksey/broker/src/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	addr := ":3000"
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("can't listen addr", zap.String("addr", addr), zap.Error(err))
	}

	srv := grpc.NewServer(
		grpc_middleware.WithStreamServerChain(
			grpc_zap.StreamServerInterceptor(logger),
		),
	)
	brk := broker.NewBroker(logger)
	pb.RegisterMessageBrokerServer(srv, transport.NewTransport(brk))

	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)

	go func() {
		logger.Debug("ready to serve")
		if err := srv.Serve(lsn); err != nil {
			log.Fatalf("server error: %v", err)
		}
		logger.Debug("stop to serve")
	}()

	logger.Debug("waiting for signal")
	s := <-sign
	logger.Debug("received signal", zap.String("signal", s.String()))

	done := make(chan struct{}, 1)
	go func() {
		srv.GracefulStop()
		done <- struct{}{}
	}()

	const gracefulTimeout = time.Second * 5
	select {
	case <-time.After(gracefulTimeout):
		logger.Debug("timed out")
	case <-done:
		logger.Debug("stopped gracefully")
	}
}
