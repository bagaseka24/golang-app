package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type HealthHandler struct {
	DB     *pgxpool.Pool
	S3     *minio.Client
	Bucket string
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "OK"})
}

// GET /health/db
func (h *HealthHandler) HealthDB(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := h.DB.Ping(ctx)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ERROR",
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  "OK",
		Message: "Database connected",
	})
}

// GET /health/s3
func (h *HealthHandler) HealthS3(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	exists, err := h.S3.BucketExists(ctx, h.Bucket)

	if err != nil || !exists {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ERROR",
			Message: "Bucket not accessible",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  "OK",
		Message: "S3 storage connected",
	})
}