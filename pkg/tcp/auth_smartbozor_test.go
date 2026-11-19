package tcp

import (
	"net/url"
	"testing"
)

func TestDisableBasic(t *testing.T) {
	DisableBasic = true
	defer func() { DisableBasic = false }()
	res := &Response{Header: map[string][]string{"Www-Authenticate": {`Basic realm="x"`}}}
	a := NewAuth(url.UserPassword("admin", "secret"))
	if a.Read(res) {
		t.Fatal("Basic rad etilishi kerak")
	}
	res.Header["Www-Authenticate"] = []string{`Basic realm="x"`, `Digest realm="r", nonce="n"`}
	if !a.Read(res) || a.Method != AuthDigest {
		t.Fatalf("Digest tanlanishi kerak: %v", a.Method)
	}
}
