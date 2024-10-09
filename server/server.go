package server

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/nicjohnson145/minibin/config"
	"github.com/rs/zerolog"
	pb "github.com/nicjohnson145/minibin/protobuf"
)

type ServerConfig struct {
	Logger     zerolog.Logger
	TemplateFS fs.FS
	FeatureSet *config.FeatureSet
}

func NewServer(conf ServerConfig) *Server {
	return &Server{
		log:        conf.Logger,
		templateFS: conf.TemplateFS,
		featureSet: conf.FeatureSet,
	}
}

type Server struct {
	pb.UnimplementedMinibinServiceServer

	log        zerolog.Logger
	templateFS fs.FS
	featureSet *config.FeatureSet
}

func (s *Server) returnError(w http.ResponseWriter, err error, context string) {
	s.log.Err(err).Msg(context)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(
			s.templateFS,
			filepath.Join("layout", "base.html"),
			filepath.Join("pages", "upload.html"),
		)
		if err != nil {
			s.returnError(w, err, "error parsing templates for home route")
			return
		}

		type homeData struct {
			Features *config.FeatureSet
		}

		data := homeData{
			Features: s.featureSet,
		}
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
			s.returnError(w, err, "error executing templates for home route")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}
