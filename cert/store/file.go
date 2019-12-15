package store

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ctrlaltdel121/cert-server/cert"
)

// FileStore implements the CertStorer interface to store certs on disk
type FileStore struct {
	StorageDir string
}

func (f FileStore) Write(cert *cert.Cert) error {
	certFile, err := os.Create(f.StorageDir + fmt.Sprintf("/%d.crt", cert.Serial))
	if err != nil {
		return err
	}
	defer certFile.Close()

	_, err = certFile.Write(cert.Cert)
	if err != nil {
		return err
	}

	// keyFile, err := os.Create(f.StorageDir + fmt.Sprintf("/%d.key", cert.Serial))
	// if err != nil {
	// 	return err
	// }
	// defer keyFile.Close()

	// _, err = keyFile.Write(cert.Key)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (f FileStore) Read(serial int64) (*cert.Cert, error) {
	var certificate cert.Cert
	var err error

	certificate.Cert, err = ioutil.ReadFile(f.StorageDir + fmt.Sprintf("/%d.crt", serial))
	if err != nil {
		return nil, err
	}
	// certificate.Key, err = ioutil.ReadFile(f.StorageDir + fmt.Sprintf("/%d.key", serial))
	// if err != nil {
	// 	return nil, err
	// }
	certificate.Serial = serial
	return &certificate, nil
}

// Delete removes the file for the given certificate serial
func (f FileStore) Delete(serial int64) error {
	return os.Remove(f.StorageDir + fmt.Sprintf("/%d.crt", serial))
}
