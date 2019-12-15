package cert

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertGenerate(t *testing.T) {
	twoD := time.Now().Add(time.Hour * 24 * 2)
	in := CertInput{
		Names:        []string{"name1", "name2"},
		ValidFrom:    nil,
		ValidTo:      &twoD,
		IsCA:         false,
		OrgName:      "Acme",
		CountryName:  "USA",
		StateName:    "NY",
		LocalityName: "NYC",
		OrgUnit:      "Acme Devops",
		CommonName:   "anotherName",
		EmailAddr:    "test@example.com",
	}

	c, err := in.Generate()
	require.NoError(t, err)

	block, _ := pem.Decode(c.Cert)
	require.NotNil(t, block)
	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	assert.EqualValues(t, []string{"name1", "name2"}, cert.DNSNames)
	assert.WithinDuration(t, time.Now().Add(time.Hour*24*2), cert.NotAfter, time.Second*2)
	assert.Equal(t, "Acme", cert.Subject.Organization[0])
	assert.Equal(t, "USA", cert.Subject.Country[0])
	assert.Equal(t, "NY", cert.Subject.Province[0])
	assert.Equal(t, "NYC", cert.Subject.Locality[0])
	assert.Equal(t, "Acme Devops", cert.Subject.OrganizationalUnit[0])
	assert.Equal(t, "anotherName", cert.Subject.CommonName)
	assert.Equal(t, "test@example.com", cert.EmailAddresses[0])
}
