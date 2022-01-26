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
	Database       string        `env:"DATABASE"`
	Host           string        `env:"DB_HOST,required"`
	Port           int           `env:"DB_PORT" envDefault:"5432"`
	User           string        `env:"DB_USER,required"`
	Pass           string        `env:"DB_PASS,required"`
	RepeatInterval time.Duration `env:"REPEAT_INTERVAL" envDefault:"0s"`
	Ssl            string        `env:"SSLMODE" envDefault:"verify-ca"`
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

	log.Info().
		Str("Database Host", host).
		Int("Database Port", cfg.Port).
		Str("Database User", user).
		Str("Database Password", pass).
		Str("Database Name", cfg.Database).
		Str("Encryption", cfg.Ssl).
		Msg("Starting Postgres Check")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s", host, cfg.Port, user, pass)

	if cfg.Database != "" {
		psqlInfo = psqlInfo + " dbname=" + cfg.Database
	}

	if cfg.Ssl != "" {
		psqlInfo = psqlInfo + " sslmode=" + cfg.Ssl
	}

	for {

		err := connectToDatabase(psqlInfo)
		if err != nil {
			return err
		}
		if cfg.RepeatInterval == time.Duration(0) {
			return nil
		}

		time.Sleep(cfg.RepeatInterval)
	}

	return nil
}

func connectToDatabase(psqlInfo string) error {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}
	log.Info().
		Msg("Database ping success")
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
