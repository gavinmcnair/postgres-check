package main

import (
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
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

	user, err := returnFileContentsOrPassword(cfg.User)
	if err != nil {
		return err
	}
	pass, err := returnFileContentsOrPassword(cfg.Pass)
	if err != nil {
		return err
	}
	host, err := returnFileContentsOrPassword(cfg.Host)
	if err != nil {
		return err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s", host, cfg.Port, user, pass)

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

func returnFileContentsOrPassword(potentialPassword string) (string, error) {
	if _, err := os.Stat(potentialPassword); err == nil {
		passwordFromFile, err := readFileAndReturnContents(potentialPassword)
		if err != nil {
			return "", err
		}
		return passwordFromFile, nil
	}

	return potentialPassword, nil
}

func readFileAndReturnContents(filename string) (string, error) {
	filebytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(filebytes), nil
}
