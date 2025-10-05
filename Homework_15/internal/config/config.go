package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
)

type Config struct {
	Addr string `json:"addr"`

	PGConn string `json:"pg_conn"`
}

func defaultConfig() Config {
	return Config{
		Addr:   ":8080",
		PGConn: "",
	}
}

func Load() (Config, error) {
	cfg := defaultConfig()

	var (
		flagAddr  = flag.String("addr", "", "HTTP listen address, e.g. :8080 or 0.0.0.0:8080")
		flagPG    = flag.String("pg", "", "PostgreSQL connection string (same as PG_CONN)")
		flagConf1 = flag.String("config", "", "Path to JSON config file")
		flagConf2 = flag.String("c", "", "Path to JSON config file (shorthand)")
	)

	flag.Parse()

	configPath := ""
	if *flagConf1 != "" {
		configPath = *flagConf1
	} else if *flagConf2 != "" {
		configPath = *flagConf2
	} else if v := os.Getenv("CONFIG"); v != "" {
		configPath = v
	}

	if configPath != "" {
		if err := loadFromJSON(configPath, &cfg); err != nil {
			return cfg, err
		}
	}

	if v := os.Getenv("PG_CONN"); v != "" {
		cfg.PGConn = v
	}

	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}

	if *flagAddr != "" {
		cfg.Addr = *flagAddr
	}
	if *flagPG != "" {
		cfg.PGConn = *flagPG
	}

	return cfg, nil
}

func loadFromJSON(path string, cfg *Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(cfg); err != nil {
		return err
	}

	if dec.More() {
		return errors.New("extra data after JSON object in config")
	}
	return nil
}
