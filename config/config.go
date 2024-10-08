package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	LogDebug = "log.debug"
	LogTrace = "log.trace"
	LogJson  = "log.json"

	FeatureAllowFileUploads       = "feature.allow_file_uploads"
	FeaturePasswordProtectUploads = "feature.password_protect_uploads"

	GrpcPort = "grpc_port"
	HttpPort = "http_port"
)

const (
	DefaultLogDebug = false
	DefaultLogTrace = false
	DefaultLogJson  = false

	DefaultFeatureAllowFileUploads       = true
	DefaultFeaturePasswordProtectUploads = false

	DefaultGrpcPort = "9090"
	DefaultHttpPort = "3050"
)

func InitializeServerConfig() {
	viper.SetDefault(LogDebug, DefaultLogDebug)
	viper.SetDefault(LogTrace, DefaultLogTrace)
	viper.SetDefault(LogJson, DefaultLogJson)

	viper.SetDefault(GrpcPort, DefaultGrpcPort)
	viper.SetDefault(HttpPort, DefaultHttpPort)

	viper.SetDefault(FeatureAllowFileUploads, DefaultFeatureAllowFileUploads)
	viper.SetDefault(FeaturePasswordProtectUploads, DefaultFeaturePasswordProtectUploads)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

type FeatureSet struct {
	AllowFileUploads       bool
	PasswordProtectUploads bool
}

func ConstructFeatureSetFromEnv() *FeatureSet {
	return &FeatureSet{
		AllowFileUploads:       viper.GetBool(FeatureAllowFileUploads),
		PasswordProtectUploads: viper.GetBool(FeaturePasswordProtectUploads),
	}
}
