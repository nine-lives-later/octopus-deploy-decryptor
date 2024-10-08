package decryptor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	examplePassword = "pass!w0rd"
)

func TestDecryptString(t *testing.T) {
	const masterKey = "dEa8P4mxooIcSeakw8vvJw=="
	const encrypted = "B4OBO6VN2oNq1BM+PETW6Y27hkClRzFsiubWQLjYd7ks9rjz0/JsppRgmGbg59+f|Ka0Ws0lTnUnqEelkKjs3yQ=="
	const decrypted = "AlwaysLookOnTheBrightSideOfLife_and_also_42"

	key, err := KeyFromMasterKey(masterKey)
	if err != nil {
		t.Fatal(err)
	}

	v, err := DecryptString(key, encrypted)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, decrypted, v)
}

func TestDecryptStringPassword_Hello(t *testing.T) {
	const encrypted = "GtoxXUnzUnBfMPNSX3bm2w==|lS4OTkwv/0HOKgUcyGc8Cg=="
	const decrypted = "Cephalopod"

	key, err := KeyFromPassword(examplePassword)
	if err != nil {
		t.Fatal(err)
	}

	v, err := DecryptString(key, encrypted)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, decrypted, v)
}

func TestDecryptStringPassword_SensitiveValue(t *testing.T) {
	const encrypted = "KSqjoSHipI2x1CrQTiFArw==|gWFgJCLGJF9OeK8YwJK0nA=="
	const decrypted = "hellow_world"

	key, err := KeyFromPassword(examplePassword)
	if err != nil {
		t.Fatal(err)
	}

	v, err := DecryptString(key, encrypted)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, decrypted, v)
}

func TestDecryptStringPassword_ScopedSensitiveValue(t *testing.T) {
	const encrypted = "0ZnH61OawGaC0VI6qZKjyw==|0jU803niVLq2F+kWj4Fg2Q=="
	const decrypted = "other_value"

	key, err := KeyFromPassword(examplePassword)
	if err != nil {
		t.Fatal(err)
	}

	v, err := DecryptString(key, encrypted)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, decrypted, v)
}
