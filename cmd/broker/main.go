package main

import (
	"flag"
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
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	addr := flag.String("listen-addr", ":8000", "Listen address")
	debug := flag.Bool("debug", false, "Run with debug logs")
	pprof := flag.Bool("pprof", false, "Run with pprof")
	flag.Parse()

	config := zap.NewProductionConfig()
	config.Level.SetLevel(zap.ErrorLevel)
	if *debug {
		config = zap.NewDevelopmentConfig()
	}
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	if *pprof {
		go func() {
			http.ListenAndServe(":8080", nil)
		}()
	}

	lsn, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Fatal("can't listen addr", zap.String("addr", *addr), zap.Error(err))
	}
	logger.Info("listen on", zap.String("addr", *addr))

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
