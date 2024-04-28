package server

import (
	"context"
	"encoding/base64"
	"html/template"
	"io/ioutil"
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
	zerolog.DefaultContextLogger = &log

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
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	h := c.
		Append(httpin.NewInput(api.CollageRequest{})).
		ThenFunc(api.Collage)

	router := http.NewServeMux()
	router.Handle("GET /collage", h)

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

	router.HandleFunc("/ar", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("cube").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Three.js Cube Example</title>
    <style>
        body { margin: 0; }
        canvas { display: block; }
    </style>
</head>
<body>
    <script type="importmap">
    {
        "imports": {
            "three": "https://unpkg.com/three@0.156.0/build/three.module.js",
            "orbitcontrols": "https://unpkg.com/three@0.164.1/examples/jsm/controls/OrbitControls.js"
        }
    }
    </script>
    <script type="module">
        import * as THREE from 'three';
        import { OrbitControls } from 'orbitcontrols';

        let camera, scene, renderer;
        let mesh;

        init();
        animate();

        function init() {
            camera = new THREE.PerspectiveCamera( 70, window.innerWidth / window.innerHeight, 0.1, 100 );
            camera.position.z = 2;
            scene = new THREE.Scene();

            const loader = new THREE.TextureLoader();
            const materials = [
                new THREE.MeshBasicMaterial({ map: loader.load('{{.Face1}}') }),
                new THREE.MeshBasicMaterial({ map: loader.load('{{.Face2}}') }),
                new THREE.MeshBasicMaterial({ map: loader.load('{{.Face1}}') }),
                new THREE.MeshBasicMaterial({ map: loader.load('{{.Face1}}') }),
                new THREE.MeshBasicMaterial({ map: loader.load('{{.Face1}}') }),
                new THREE.MeshBasicMaterial({ map: loader.load('{{.Face1}}') })
            ];

            const geometry = new THREE.BoxGeometry();
            mesh = new THREE.Mesh(geometry, materials);
            scene.add(mesh);

            renderer = new THREE.WebGLRenderer({ antialias: true });
            renderer.setPixelRatio(window.devicePixelRatio);
            renderer.setSize(window.innerWidth, window.innerHeight);
            document.body.appendChild(renderer.domElement);
            window.addEventListener('resize', onWindowResize);

            const controls = new OrbitControls(camera, renderer.domElement);
            controls.enableDamping = true;
            controls.dampingFactor = 0.25;
            controls.enableZoom = true;
        }

        function onWindowResize() {
            camera.aspect = window.innerWidth / window.innerHeight;
            camera.updateProjectionMatrix();
            renderer.setSize(window.innerWidth, window.innerHeight);
        }

        function animate() {
            requestAnimationFrame(animate);
            mesh.rotation.x += 0.0010;
            mesh.rotation.y += 0.0010;
            renderer.render(scene, camera);
        }
    </script>
</body>
</html>
`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Create an instance of TemplateData with the desired paths
		imagebase64, _ := tobase64("assets/songstitch_logo_dark.png")
		data := TemplateData{
			Face1: "data:image/png;base64," + imagebase64,
			Face2: "data:image/png;base64, " + imagebase64,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Minute,
	}

	http.DefaultClient.Timeout = 10 * time.Second
	spotify.InitSpotifyClient(context.Background())

	log.Info().Msg("Starting server...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}

type TemplateData struct {
	Face1 string
	Face2 string
}

func tobase64(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
