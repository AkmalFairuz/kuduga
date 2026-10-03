package provider

import (
	"io"
)

type StorageProvider interface {
	Put(bucket, path string, buf io.Reader) error
	Get(bucket, path string) (io.Reader, error)
	PublicEndpoint() string

	// List TODO: return detailed object
	List(bucket, prefix string) ([]string, error)
}
