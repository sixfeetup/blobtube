package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/sixfeetup/blobtube/internal/api"
	"github.com/sixfeetup/blobtube/internal/config"
)

type runOptions struct {
	signals []os.Signal
}

type RunOption func(*runOptions)

func WithSignals(sig ...os.Signal) RunOption {
	return func(o *runOptions) {
		o.signals = append([]os.Signal(nil), sig...)
	}
}

func Run(ctx context.Context, cfg config.Config, opts ...RunOption) error {
	o := runOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	if len(o.signals) == 0 {
		o.signals = []os.Signal{os.Interrupt}
	}

	ctx, stop := signal.NotifyContext(ctx, o.signals...)
	defer stop()

	h, err := api.NewHandler(cfg)
	if err != nil {
		return err
	}

	httpsSrv := &http.Server{
		Addr:              cfg.HTTPSAddr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12},
	}

	httpRedirectSrv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           redirectToHTTPS(cfg),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 2)

	go func() {
		log.Info().Str("addr", cfg.HTTPSAddr).Msg("https server starting")
		errCh <- httpsSrv.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
	}()

	go func() {
		log.Info().Str("addr", cfg.HTTPAddr).Msg("http redirect server starting")
		errCh <- httpRedirectSrv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = httpRedirectSrv.Shutdown(shutdownCtx)
		_ = httpsSrv.Shutdown(shutdownCtx)

		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return err
	}
}

func redirectToHTTPS(cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if host == "" {
			host = "localhost"
		}

		hostname := host
		httpsPort := strings.TrimPrefix(cfg.HTTPSAddr, ":")
		if h, _, splitErr := net.SplitHostPort(host); splitErr == nil {
			hostname = h
		}
		if strings.Contains(hostname, ":") && !strings.HasPrefix(hostname, "[") {
			hostname = "[" + hostname + "]"
		}

		// Always redirect to the HTTPS listener port.
		location := fmt.Sprintf("https://%s:%s%s", hostname, httpsPort, r.URL.RequestURI())
		http.Redirect(w, r, location, http.StatusTemporaryRedirect)
	})
}
