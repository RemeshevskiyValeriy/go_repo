package config

import (
	"os"
	
	"github.com/joho/godotenv"
)

type Config struct {
	Addr     string
	CertFile string
	KeyFile  string
	DSN      string
}

func New() Config {
	_ = godotenv.Load()

	ADDR := os.Getenv("ADDR")
	
	CERT_FILE := os.Getenv("CERT_FILE")
	KEY_FILE := os.Getenv("KEY_FILE")
	
	DSN := os.Getenv("DATABASE_URL")
	
	return Config{
		Addr:     ADDR,
		CertFile: CERT_FILE,
		KeyFile:  KEY_FILE,
		DSN:      DSN,
	}
}
