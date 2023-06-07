package main

import (
	"github.com/SongStitch/song-stitch/internal/server"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("Error loading .env file")
	}

	server.RunServer()
}
