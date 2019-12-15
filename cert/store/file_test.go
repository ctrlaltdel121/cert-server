package store

import (
	"os"
	"testing"

	"github.com/ctrlaltdel121/cert-server/cert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileOps(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)
	fs := FileStore{dir}

	crt := cert.Cert{
		Cert:   []byte("dummydata"),
		Key:    []byte("dummydata"),
		Serial: 1,
	}

	// cleanup in case test fails
	defer func() {
		os.Remove(dir + "/1.crt")
	}()

	err = fs.Write(&crt)
	require.NoError(t, err)

	// confirm the cert was written
	_, err = os.Stat(dir + "/1.crt")
	require.NoError(t, err)

	newCrt, err := fs.Read(1)
	require.NoError(t, err)

	assert.Equal(t, []byte("dummydata"), newCrt.Cert)
	assert.Equal(t, int64(1), newCrt.Serial)

	err = fs.Delete(1)
	require.NoError(t, err)

	_, err = os.Stat(dir + "/1.crt")
	require.True(t, os.IsNotExist(err))
}
