package http

import (
	"context"
	"net"
	"net/http"
)

// Server represents the http server, all config concerning http is done in here
type Server struct {
	ln     net.Listener
	server *http.Server
	router http.Handler
}

func NewServer(addr string, router http.Handler) Server {
	return Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           router,
			TLSConfig:         nil,
			ReadTimeout:       0,
			ReadHeaderTimeout: 0,
			WriteTimeout:      0,
			IdleTimeout:       0,
			MaxHeaderBytes:    0,
			TLSNextProto:      nil,
			ConnState:         nil,
			ErrorLog:          nil,
			BaseContext:       nil,
			ConnContext:       nil,
		},
	}
}

// Run runs a http server listening on the port configured.
// this ends when ListenAndServe receives an error, or the context is cancelled
func (s Server) Run(ctx context.Context) (err error) {
	if s.ln, err = net.Listen("tcp", s.server.Addr); err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		errCh <- s.server.Serve(s.ln)
	}()
	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

// Shutdown runs the sver shutdown method which gracefully finished off any current requests,
// and prevents new ones from coming in
func (s Server) Shutdown(ctx context.Context) error {
	if err := s.ln.Close(); err != nil {
		return err
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
