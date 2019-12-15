package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math"
	"math/big"
	"strconv"
	"time"
)

// CertInput is filled in by users of this package, and then the users
// can call validate on the struct to get back the self-signed cert
type CertInput struct {
	Names        []string   `json:"names"`
	ValidFrom    *time.Time `json:"valid_from"`
	ValidTo      *time.Time `json:"valid_to"`
	IsCA         bool       `json:"is_ca"`
	OrgName      string     `json:"organization_name"`
	CountryName  string     `json:"country_name"`
	StateName    string     `json:"state_name"`
	LocalityName string     `json:"locality_name"`
	OrgUnit      string     `json:"organizational_unit"`
	CommonName   string     `json:"common_name"`
	EmailAddr    string     `json:"email_address"`
}

// Cert contains PEM data for a cert and it's private key
type Cert struct {
	Cert   []byte
	Key    []byte
	Serial int64
}

// Generate returns two byte slices, the PEM cert and PEM key, and an error if generation fails.
func (c *CertInput) Generate() (*Cert, error) {
	// automatic default values
	if c.ValidFrom == nil {
		now := time.Now()
		c.ValidFrom = &now
	}
	if c.ValidTo == nil {
		f := time.Now().Add(30 * time.Hour * 24)
		c.ValidTo = &f
	}
	if c.OrgName == "" {
		c.OrgName = "Acme Inc"
	}

	// if they don't specify a common name, use the first DNS name specified.
	// modern SSL certs use DNS SAN section instead of relying on CN.
	if c.CommonName == "" {
		c.CommonName = c.Names[0]
	}

	// Generate private key for cert - could make bit size configurable
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// best practice for CAs assigning serials is to generate one randomly.
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	// build certificate
	x509Cert := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   c.CommonName,
			SerialNumber: strconv.FormatInt(serialNumber.Int64(), 10),
		},
		NotBefore:             *c.ValidFrom,
		NotAfter:              *c.ValidTo,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              c.Names,
	}

	if c.IsCA {
		x509Cert.IsCA = true
		x509Cert.KeyUsage |= x509.KeyUsageCertSign
	}

	// add optional values to subject if they aren't empty
	if c.CountryName != "" {
		x509Cert.Subject.Country = []string{c.CountryName}
	}
	if c.OrgName != "" {
		x509Cert.Subject.Organization = []string{c.OrgName}
	}
	if c.OrgUnit != "" {
		x509Cert.Subject.OrganizationalUnit = []string{c.OrgUnit}
	}
	if c.LocalityName != "" {
		x509Cert.Subject.Locality = []string{c.LocalityName}
	}
	if c.StateName != "" {
		x509Cert.Subject.Province = []string{c.StateName}
	}
	if c.EmailAddr != "" {
		x509Cert.EmailAddresses = []string{c.EmailAddr}
	}

	// Get certificate DER bytes
	derBytes, err := x509.CreateCertificate(rand.Reader, &x509Cert, &x509Cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	// encodeToMemory takes the bytes, base64 encodes them, and wraps them in the PEM block (BEGIN/END header/footer)
	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return &Cert{cert, key, serialNumber.Int64()}, nil
}
