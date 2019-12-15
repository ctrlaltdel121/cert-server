package store

import "github.com/ctrlaltdel121/cert-server/cert"

// This is an example of how you could use the interface
// to implement an S3 store instead of on local disk.

// S3Store implements the CertStorer interface to store certs on S3
type S3Store struct {
	BucketName         string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
}

func (f S3Store) Write(cert *cert.Cert) error {
	// TODO
	return nil
}

func (f S3Store) Read(serial int64) (*cert.Cert, error) {
	// TODO
	return nil, nil
}

func (f S3Store) Delete(serial int64) error {
	// TODO
	return nil
}
