package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/networkteam/slogutils"

	"github.com/esdete2/envelopr/handler"
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
	s.router.Route("/_template", func(r chi.Router) {
		r.Get("/*", s.handleRawTemplate())
	})
	s.router.Route("/_events", func(r chi.Router) {
		r.Get("/", s.broker.ServeHTTP)
	})
	s.router.Get("/*", s.handleTemplate())
}

func (s *Server) Serve(addr string) error {
	slog.With("address", "http://"+addr).Info("Server started")
	//nolint: gosec
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown failed", slogutils.Err(err))
		}
	}()

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
				).Info("request completed")
			}()

			next.ServeHTTP(responseWriter, r)
		})
	}
}
