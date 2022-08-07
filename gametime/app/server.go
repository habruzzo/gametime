package app

import (
	"context"
	"gametime"
	"gametime/config"
	"net"

	"net/http"
)

type Router interface {
	http.Handler
}

func NewServer(log gametime.Logger, router Router, cfg *config.Config) *Server {

	return &Server{

		log: log,

		router: router,

		addr: cfg.Address(),
	}

}

type Server struct {
	router Router

	addr string

	log gametime.Logger
}

func (s *Server) Name() string {

	return "api server"

}

func (s *Server) Start(ctx context.Context) error {

	listener, err := net.Listen("tcp", s.addr)

	if err != nil {

		return gametime.Error{

			Actual: err,

			Type: gametime.Unrecoverable,
		}

	}

	httpServer := &http.Server{

		Handler: s.router,

		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	// We declare the listener here so if we have issues listening we can fail to start

	s.log.Info(ctx, "Listening on port ", listener.Addr())

	go func() {

		defer httpServer.Close()

		<-ctx.Done()

	}()

	return gametime.Error{

		Actual: httpServer.Serve(listener),

		Type: gametime.Unrecoverable,
	}

}
