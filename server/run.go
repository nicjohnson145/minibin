package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"

	"github.com/nicjohnson145/minibin/config"
	"github.com/spf13/viper"
)

//go:embed ui
var uiFS embed.FS

func Run() error {
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

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assetFS))))
	mux.Handle("/", server.Home())

	lis, err := net.Listen("tcp4", ":"+viper.GetString(config.Port))
	if err != nil {
		return fmt.Errorf("error starting listener: %w", err)
	}

	srv := &http.Server{
		Handler: mux,
	}
	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}
