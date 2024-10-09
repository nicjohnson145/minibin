package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicjohnson145/minibin/config"
	pb "github.com/nicjohnson145/minibin/protobuf"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

//go:embed ui
var uiFS embed.FS

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.InitializeServerConfig()
	log := config.InitLogger()


	assetFS, err := fs.Sub(uiFS, "ui/assets")
	if err != nil {
		return fmt.Errorf("error building asset FS: %w", err)
	}

	templateFS, err := fs.Sub(uiFS, "ui/templates")
	if err != nil {
		return fmt.Errorf("error building template FS: %w", err)
	}

	server := NewServer(ServerConfig{
		Logger:     config.WithComponent(log, "server"),
		TemplateFS: templateFS,
		FeatureSet: config.ConstructFeatureSetFromEnv(),
	})

	grpcServer := grpc.NewServer()
	pb.RegisterMinibinServiceServer(grpcServer, server)
	reflection.Register(grpcServer)


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	grpcPort := viper.GetString(config.GrpcPort)
	grpcLis, err := net.Listen("tcp4", ":" + grpcPort)
	if err != nil {
		return fmt.Errorf("error starting grpc listener: %w", err)
	}

	httpPort := viper.GetString(config.HttpPort)
	httpLis, err := net.Listen("tcp4", ":" + httpPort)
	if err != nil {
		return fmt.Errorf("error starting http listener: %w", err)
	}

	grpcMux := runtime.NewServeMux()
	if err := pb.RegisterMinibinServiceHandlerFromEndpoint(ctx, grpcMux, "localhost:" + grpcPort, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
		return fmt.Errorf("error registering gateway: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", grpcMux)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assetFS))))
	mux.Handle("/", server.Home())

	httpServer := http.Server{
		Handler: mux,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		stop := func() {
			grpcServer.GracefulStop()
			httpServer.Shutdown(ctx)
			wg.Done()
		}
		select {
		case s := <- sigChan:
			log.Info().Msgf("got signal %v, attempting shutdown", s)
			cancel()
			stop()
		case <-ctx.Done():
			log.Info().Msg("context cancelled, attempting graceful shutdown")
			stop()
		}
	}()

	go func() {
		log.Info().Msgf("starting http server on :%v", httpPort)
		if err := httpServer.Serve(httpLis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancel()
			log.Err(err).Msg("error serving http")
		}
		log.Info().Msg("http server goroutine exiting")
	}()

	go func() {
		log.Info().Msgf("starting grpc server on :%v", grpcPort)
		if err := grpcServer.Serve(grpcLis); err != nil {
			cancel()
			log.Err(err).Msg("error serving grpc")
		}
		log.Info().Msg("grpc server goroutine exiting")
	}()

	wg.Wait()
	return nil
}
