package srv

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ctrlaltdel121/cert-server/cert"
	"github.com/gavv/httpexpect"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerBasic(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	defer func() {
		os.Remove(dir)
	}()
	os.Setenv("STORAGE_DIR", dir)
	srv := NewServer("file")

	tst := apiTester(t, srv.createRouter())

	// test creation of a cert, minimal example
	resp := tst.POST("/certificates").WithBytes([]byte(`{"names":["name1", "name2"]}`)).Expect()
	resp.Status(201)
	bodyStr := resp.Body().Raw()
	certResp := cert.Cert{}
	err = json.Unmarshal([]byte(bodyStr), &certResp)
	require.NoError(t, err)

	block, _ := pem.Decode(certResp.Cert)
	require.NotNil(t, block)

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	assert.EqualValues(t, []string{"name1", "name2"}, parsedCert.DNSNames)
	assert.WithinDuration(t, time.Now().Add(30*time.Hour*24), parsedCert.NotAfter, time.Second*2)
	assert.Equal(t, "Acme Inc", parsedCert.Subject.Organization[0])
	assert.Equal(t, "name1", parsedCert.Subject.CommonName)

	// test that the returned cert/key loads as a valid TLS keypair
	_, err = tls.X509KeyPair(certResp.Cert, certResp.Key)
	require.NoError(t, err)

	// test that we can use the serial to get the certificate back
	resp = tst.GET("/certificates/" + strconv.Itoa(int(certResp.Serial))).Expect()
	resp.Status(200)
	bodyStr = resp.Body().Raw()
	certResp = cert.Cert{}
	err = json.Unmarshal([]byte(bodyStr), &certResp)
	require.NoError(t, err)

	block, _ = pem.Decode(certResp.Cert)
	require.NotNil(t, block)

	parsedCert, err = x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	assert.EqualValues(t, []string{"name1", "name2"}, parsedCert.DNSNames)
	assert.WithinDuration(t, time.Now().Add(30*time.Hour*24), parsedCert.NotAfter, time.Second*2)
	assert.Equal(t, "Acme Inc", parsedCert.Subject.Organization[0])
	assert.Equal(t, "name1", parsedCert.Subject.CommonName)

	// no key after first creation
	assert.Empty(t, certResp.Key)

	// able to delete the cert, no error
	resp = tst.DELETE("/certificates/" + strconv.Itoa(int(certResp.Serial))).Expect()
	resp.Status(202)

	// cert not found
	resp = tst.GET("/certificates/" + strconv.Itoa(int(certResp.Serial))).Expect()
	resp.Status(404)

}

// Test a fully filled in cert
func TestServerFullCert(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	defer func() {
		os.Remove(dir)
	}()
	os.Setenv("STORAGE_DIR", dir)
	srv := NewServer("file")

	tst := apiTester(t, srv.createRouter())

	// test creation of a cert, full example
	resp := tst.POST("/certificates").WithBytes([]byte(`{
		"names":["name1", "name2", "name3"],
		"common_name":"testname",
		"valid_to":"2019-10-20T10:10:10Z",
		"is_ca": true,
		"organization_name":"Acme",
		"organizational_unit":"testers",
		"country_name":"USA",
		"state_name":"NY",
		"locality_name":"NYC",
		"email_address":"test@example.com"
	}`)).Expect()
	resp.Status(201)
	bodyStr := resp.Body().Raw()
	certResp := cert.Cert{}
	err = json.Unmarshal([]byte(bodyStr), &certResp)
	require.NoError(t, err)

	block, _ := pem.Decode(certResp.Cert)
	require.NotNil(t, block)

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	tm, err := time.Parse(time.RFC3339, "2019-10-20T10:10:10Z")
	require.NoError(t, err)
	assert.EqualValues(t, []string{"name1", "name2", "name3"}, parsedCert.DNSNames)
	assert.Equal(t, tm, parsedCert.NotAfter)
	assert.Equal(t, "Acme", parsedCert.Subject.Organization[0])
	assert.Equal(t, "testers", parsedCert.Subject.OrganizationalUnit[0])
	assert.Equal(t, "USA", parsedCert.Subject.Country[0])
	assert.Equal(t, "NY", parsedCert.Subject.Province[0])
	assert.Equal(t, "NYC", parsedCert.Subject.Locality[0])
	assert.Equal(t, "test@example.com", parsedCert.EmailAddresses[0])
	assert.Equal(t, "testname", parsedCert.Subject.CommonName)
	assert.True(t, parsedCert.IsCA)

	// test that the returned cert/key loads as a valid TLS keypair
	_, err = tls.X509KeyPair(certResp.Cert, certResp.Key)
	require.NoError(t, err)

	resp = tst.DELETE("/certificates/" + strconv.Itoa(int(certResp.Serial))).Expect()
	resp.Status(202)

}

func apiTester(t *testing.T, r *mux.Router) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://example.com", // doesn't matter
		Client: &http.Client{
			Transport: httpexpect.NewBinder(r),
		},
		Reporter: httpexpect.NewRequireReporter(t),
		Printers: []httpexpect.Printer{httpexpect.NewDebugPrinter(t, true)},
	})
}
