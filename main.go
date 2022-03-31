package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	ListenPort     int           `env:"LISTEN_PORT" envDefault:"8080"`
	RepeatInterval time.Duration `env:"REPEAT_INTERVAL" envDefault:"15s"`
}

func main() {
	zerolog.DurationFieldUnit = time.Second
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}
	log.Info().Msg("Gracefully exiting")
}

func run() error {
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		return err
	}

	createPrometheusEndpoint(cfg).ListenAndServe()
	return nil
}

func createPrometheusEndpoint(cfg config) *http.Server {


	log.Info().
 		Int("Metrics Port", cfg.ListenPort).
		Str("Path", "/metrics").
		Msg("Prometheus endpoint service started.")

	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("RequestFrom", r.RemoteAddr).
			Msg("Metrics Collection has occured")
		promhttp.Handler().ServeHTTP(w, r)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})

	return &http.Server{
		Handler:      mux,
		Addr:         ":" + strconv.Itoa(cfg.ListenPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
