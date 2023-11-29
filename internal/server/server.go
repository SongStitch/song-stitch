package server

import (
	"net/http"
	"os"
	"time"

	"github.com/ggicci/httpin"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/SongStitch/song-stitch/internal/api"
	"github.com/SongStitch/song-stitch/internal/clients/spotify"
	"github.com/SongStitch/song-stitch/internal/config"
)

func getLogger() zerolog.Logger {
	if env, ok := os.LookupEnv("APP_ENV"); ok && env == "development" {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		return zerolog.New(output).With().Timestamp().Logger()
	} else {
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

func RunServer() {
	_ = godotenv.Load()
	log := getLogger()

	err := config.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize config")
	}

	c := alice.New()
	c = c.Append(hlog.NewHandler(log))
	c = c.Append(
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
	)
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	h := c.
		Append(httpin.NewInput(api.CollageRequest{})).
		ThenFunc(api.Collage)

	router := http.NewServeMux()
	router.Handle("/collage", h)

	// serve files from public folder
	fs := http.FileServer(http.Dir("./public"))
	router.Handle("/", fs)

	// serve robots.txt file
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/robots.txt")
	})

	// serve humans.txt file
	router.HandleFunc("/humans.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/humans.txt")
	})

	// serve support page
	router.HandleFunc("/support", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/support.html")
	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Minute,
	}

	http.DefaultClient.Timeout = 10 * time.Second
	spotify.InitSpotifyClient(log)

	log.Info().Msg("Starting server...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}
