package web

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/networkteam/slogutils"

	"github.com/esdete2/mjml-dev/handler"
)

type Server struct {
	router  *chi.Mux
	options *ServerOptions
	broker  *EventBroker
}

type ServerOptions struct {
	Output string
}

func NewServer(opts *ServerOptions) *Server {
	logger := slogutils.FromContext(context.Background())
	r := chi.NewRouter()
	broker := NewEventBroker()

	// Custom logger middleware
	r.Use(loggerMiddleware(logger))
	r.Use(chimiddleware.Recoverer)

	s := &Server{
		router:  r,
		options: opts,
		broker:  broker,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.Get("/", s.handleIndex())
	s.router.Get("/_template/*", s.handleRawTemplate())
	s.router.Get("/_events", s.broker.ServeHTTP)
	s.router.Get("/*", s.handleTemplate())
}

func (s *Server) Serve(addr string) error {
	slog.With("address", "http://"+addr).Info("Server started")
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) NotifyReload() {
	slog.Debug("Notifying reload")
	s.broker.Notify("reload")
}

func (s *Server) ReloadNotifier() handler.ReloadNotifier {
	return s
}

// loggerMiddleware creates a Chi middleware for structured logging
func loggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.With("method", r.Method).
					With("path", r.URL.Path).
					With("status", ww.Status()).
					With("size", ww.BytesWritten()).
					With("duration", chimiddleware.GetReqID(r.Context())).
					Info("request completed")
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
