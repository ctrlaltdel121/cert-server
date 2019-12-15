package store

import (
	"os"
	"strings"

	"github.com/ctrlaltdel121/cert-server/cert"
)

// CertStorer defines an interface for storing certs
type CertStorer interface {
	Write(cert *cert.Cert) error
	Read(serial int64) (*cert.Cert, error)
	Delete(serial int64) error
}

// NewFileStore gives a new CertStorer that stores files in STORAGE_DIR
func NewFileStore() CertStorer {
	return FileStore{strings.TrimSuffix(os.Getenv("STORAGE_DIR"), "/")}
}

// NewS3Store gives a new CertStorer that stores files in S3 (not fully implemented)
func NewS3Store() CertStorer {
	return S3Store{}
}
