package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	BindAddress string
	BasePath    string
}

func Load() *Configs {
	godotenv.Load(".env")

	cfg := Configs{
		BindAddress: os.Getenv("BIND_ADDR"),
		BasePath:    os.Getenv("BASE_PATH"),
	}

	return &cfg
}
