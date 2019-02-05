package probe

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/iotexproject/iotex-core/pkg/log"
)

const (
	_ready    = 1
	_notReady = 0
)

// Server is a http server for service probe.
type Server struct {
	ready            int32 // 0 is not ready, 1 is ready
	server           http.Server
	readinessHandler http.Handler
}

// Option is ued to set probe server's options.
type Option interface {
	SetOption(*Server)
}

// New creates a new probe server.
func New(port int, opts ...Option) *Server {
	s := &Server{
		ready:            _notReady,
		readinessHandler: http.HandlerFunc(successHandleFunc),
	}

	for _, opt := range opts {
		opt.SetOption(s)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/liveness", successHandleFunc)
	readiness := func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&s.ready) == _notReady {
			failureHandleFunc(w, r)
			return
		}
		s.readinessHandler.ServeHTTP(w, r)
	}
	mux.HandleFunc("/readiness", readiness)
	mux.HandleFunc("/health", readiness)

	s.server = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return s
}

// Start starts the probe server and starts returning success status on liveness endpoint.
func (s *Server) Start(_ context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.L().Info("Probe server stopped.", zap.Error(err))
		}
	}()
	return nil
}

// Ready makes the probe server starts returning status on readiness and
// health endpoint.
func (s *Server) Ready() { atomic.SwapInt32(&s.ready, _ready) }

// NotReady makes the probe server starts returning failure status on readiness and
// health endpoint.
func (s *Server) NotReady() { atomic.SwapInt32(&s.ready, _notReady) }

// Stop shutdown the probe server.
func (s *Server) Stop(ctx context.Context) error { return s.server.Shutdown(ctx) }

func successHandleFunc(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.L().Warn("Failed to send http response.", zap.Error(err))
	}
}

func failureHandleFunc(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	if _, err := w.Write([]byte("FAIL")); err != nil {
		log.L().Warn("Failed to send http response.", zap.Error(err))
	}
}