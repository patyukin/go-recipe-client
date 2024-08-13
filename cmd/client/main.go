package main

import (
	"app/internal/config"
	"app/internal/server"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"net"
	"os"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	hostClient := &fasthttp.HostClient{
		Addr: net.JoinHostPort(cfg.RemoteHost, fmt.Sprintf("%d", cfg.RemotePort)),
	}

	pipelineClient := &fasthttp.PipelineClient{
		Addr: net.JoinHostPort(cfg.RemoteHost, fmt.Sprintf("%d", cfg.RemotePort)),
	}

	srv := server.New(cfg, hostClient, pipelineClient)

	r := router.New()

	r.GET("/get-recipes", srv.GetAllRecipes)

	log.Info().Msgf("Server is running on :%d", cfg.Port)
	log.Fatal().Err(fasthttp.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r.Handler))
}
