package hwdns_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wonfen/huawei-GTM-DNS/hwdns"
)

func TestSigner(t *testing.T) {
	signer := hwdns.NewSigner("test-ak", "test-sk")

	req, err := http.NewRequest("GET", "https://dns.myhuaweicloud.com/v2/zones", nil)
	require.NoError(t, err)

	err = signer.Sign(req)
	require.NoError(t, err)

	auth := req.Header.Get("Authorization")
	assert.Contains(t, auth, "SDK-HMAC-SHA256")
	assert.Contains(t, auth, "test-ak")
	assert.NotEmpty(t, req.Header.Get("X-Sdk-Date"))
}

func TestSignerWithBody(t *testing.T) {
	signer := hwdns.NewSigner("test-ak", "test-sk")
	body := strings.NewReader(`{"status":"DISABLE"}`)
	req, err := http.NewRequest("PUT", "https://dns.myhuaweicloud.com/v2/zones/z1/recordsets/r1", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	err = signer.Sign(req)
	require.NoError(t, err)

	auth := req.Header.Get("Authorization")
	assert.Contains(t, auth, "SDK-HMAC-SHA256")
	assert.Contains(t, auth, "SignedHeaders=")
	assert.Contains(t, auth, "Signature=")
}
