package provider

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type MinioConfig struct {
	AccessKey string `env:"MINIO_ACCESS_KEY"`
	SecretKey string `env:"MINIO_SECRET_KEY"`
	Endpoint  string `env:"MINIO_ENDPOINT"`
	UseSsl    bool   `env:"MINIO_USE_SSL"`
	PublicUrl string `env:"MINIO_PUBLIC_URL"`
}

type MinioStorageProvider struct {
	minioClient *minio.Client
	publicUrl   string
}

func NewMinioStorageProvider(config MinioConfig) (*MinioStorageProvider, error) {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSsl,
	})
	if err != nil {
		return nil, err
	}
	return &MinioStorageProvider{minioClient: minioClient, publicUrl: config.PublicUrl}, nil
}

func (p MinioStorageProvider) PublicEndpoint() string {
	return p.publicUrl
}

func (p MinioStorageProvider) Put(bucket, path string, buf io.Reader) error {
	_, err := p.minioClient.PutObject(context.TODO(), bucket, path, buf, -1, minio.PutObjectOptions{})
	return err
}

func (p MinioStorageProvider) Get(bucket, path string) (io.Reader, error) {
	object, err := p.minioClient.GetObject(context.TODO(), bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (p MinioStorageProvider) List(bucket, prefix string) ([]string, error) {
	ret := make([]string, 0)
	for object := range p.minioClient.ListObjects(context.TODO(), bucket, minio.ListObjectsOptions{Prefix: prefix}) {
		ret = append(ret, object.Key)
	}
	return ret, nil
}
