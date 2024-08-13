package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	RemotePort     int
	RemoteHost     string
	HostClient     *fasthttp.HostClient
	PipelineClient *fasthttp.PipelineClient
}

func LoadConfig() (*Config, error) {
	envConfigFilePath := os.Getenv("ENV_CONFIG_FILE_PATH")
	if err := godotenv.Load(envConfigFilePath); err != nil {
		log.Fatal().Msgf("Error loading .env file")
	}

	cfg := &Config{}

	remotePort, err := strconv.Atoi(os.Getenv("REMOTE_PORT"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse REMOTE_PORT: %v", err)
	}

	cfg.RemotePort = remotePort

	port, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTTP_PORT: %v", err)
	}

	cfg.Port = port
	cfg.RemoteHost = os.Getenv("REMOTE_HOST")

	return cfg, nil
}
