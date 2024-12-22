package web

import (
	"log/slog"
	"net/http"
	"time"

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
	router := chi.NewRouter()
	broker := NewEventBroker()

	// Custom logger middleware
	router.Use(loggerMiddleware())
	router.Use(chimiddleware.Recoverer)

	srv := &Server{
		router:  router,
		options: opts,
		broker:  broker,
	}

	srv.routes()
	return srv
}

func (s *Server) routes() {
	s.router.Get("/", s.handleIndex())
	s.router.Get("/_template/*", s.handleRawTemplate())
	s.router.Get("/_events", s.broker.ServeHTTP)
	s.router.Get("/*", s.handleTemplate())
}

func (s *Server) Serve(addr string) error {
	slog.With("address", "http://"+addr).Info("Server started")
	srv := &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) NotifyReload() {
	slog.Debug("Notifying reload")
	s.broker.Notify("reload")
}

func (s *Server) ReloadNotifier() handler.ReloadNotifier {
	return s
}

// loggerMiddleware creates a Chi middleware for structured logging
func loggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseWriter := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			logger := slogutils.FromContext(r.Context())

			defer func() {
				logger.With(
					"method", r.Method,
					"path", r.URL.Path,
					"status", responseWriter.Status(),
					"size", responseWriter.BytesWritten(),
					"duration", chimiddleware.GetReqID(r.Context()),
				).Info("request completed")
			}()

			next.ServeHTTP(responseWriter, r)
		})
	}
}