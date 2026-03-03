package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type FileHandler struct {
	S3     *minio.Client
	Bucket string
}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

// POST /api/upload
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := context.Background()

	// Upload to MinIO
	info, err := h.S3.PutObject(ctx, h.Bucket, header.Filename, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	response := FileInfo{
		Name: header.Filename,
		Size: info.Size,
		URL:  fmt.Sprintf("/api/files/%s", header.Filename),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GET /api/files
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	objectCh := h.S3.ListObjects(ctx, h.Bucket, minio.ListObjectsOptions{})

	files := []FileInfo{}
	for object := range objectCh {
		if object.Err != nil {
			http.Error(w, "Failed to list files", http.StatusInternalServerError)
			return
		}

		files = append(files, FileInfo{
			Name: object.Key,
			Size: object.Size,
			URL:  fmt.Sprintf("/api/files/%s", object.Key),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// GET /api/files/:name
func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL (simple approach)
	filename := r.URL.Path[len("/api/files/"):]

	ctx := context.Background()

	object, err := h.S3.GetObject(ctx, h.Bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", stat.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size))
	io.Copy(w, object)
}