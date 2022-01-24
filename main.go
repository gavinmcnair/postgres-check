package main

import (
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type config struct {
	Database string `env:"DATABASE"`
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER,required"`
	Pass     string `env:"DB_PASS,required"`
	Ssl      string `env:"SSLMODE" envDefault:"disable"`
}

func main() {
	zerolog.DurationFieldUnit = time.Second
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}
	log.Info().Msg("Database connection successful - gracefully exiting")
}

func run() error {
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		return err
	}

	log.Info().
		Str("Database Host", cfg.Host).
		Int("Database Port", cfg.Port).
		Str("Database User", cfg.User).
		Str("Database Password", "Hidden").
		Str("Database Name", cfg.Database).
		Str("Encryption", cfg.Ssl).
		Msg("Starting Postgres Check")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s", cfg.Host, cfg.Port, cfg.User, cfg.Pass)

	if cfg.Database != "" {
		psqlInfo = psqlInfo + " dbname=" + cfg.Database
	}

	if cfg.Ssl != "" {
		psqlInfo = psqlInfo + " sslmode=" + cfg.Ssl
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}
