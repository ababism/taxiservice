package mshttp

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *Config, r *Router) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      r.router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// TODO support http server with tls

func (r *Server) Start() error {
	if r.srv.Handler == nil {
		return errors.New("no routes have been registered")
	}
	if err := r.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (r *Server) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
