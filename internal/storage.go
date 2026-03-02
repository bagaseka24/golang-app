package internal

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Client *minio.Client
	Bucket string
}

func ConnectS3(endpoint, accessKey, secretKey, bucket string) (*Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &Storage{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (s *Storage) CheckConnection(ctx context.Context) error {
	_, err := s.Client.BucketExists(ctx, s.Bucket)
	return err
}
