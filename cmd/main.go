package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"devops-app/internal"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := internal.LoadConfig()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	if cfg.DbEnabled {
		conn, err := internal.ConnectDB(cfg.DbURL)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close(context.Background())

		r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("DB connected"))
		})
	}

	var storage *internal.Storage
	if cfg.S3Enabled {
		s, err := internal.ConnectS3(
			cfg.S3Endpoint,
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			cfg.S3Bucket,
		)
		if err != nil {
			log.Println("S3 connection error:", err)
		} else {
			storage = s
		}

		r.Get("/storage/health", func(w http.ResponseWriter, r *http.Request) {
			if storage == nil {
				http.Error(w, "S3 not configured", http.StatusServiceUnavailable)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := storage.CheckConnection(ctx)
			if err != nil {
				http.Error(w, "Cannot access bucket: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write([]byte("S3 OK"))
		})
	} else {
		r.Get("/storage/health", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "S3 disabled", http.StatusServiceUnavailable)
		})
	}

	log.Println("Starting app on port", cfg.AppPort)
	http.ListenAndServe(":"+cfg.AppPort, r)
}
