package identity

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mendsley/gojwk"
	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

func MustFetchKeys(verifierKeyURI string, client HTTPClient) map[string]*rsa.PublicKey {
	for {
		var err error
		for i := 0; i < 20; i++ {
			var keys map[string]*rsa.PublicKey
			keys, err = FetchKeys(verifierKeyURI, client)
			if err == nil {
				return keys
			}

			time.Sleep(1 * time.Second)
		}

		logs.Logger.Warnf("Could not fetch keys, error %v. Continuing to try...", err)
	}
}

func FetchKeys(verifierKeyURI string, client HTTPClient) (map[string]*rsa.PublicKey, error) {
	response, err := client.Get(verifierKeyURI)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth server returned a non-200 response code, code was %d", response.StatusCode)
	}

	logs.Logger.Infof("Successfully fetched public key from %s", verifierKeyURI)

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	jwks, err := gojwk.Unmarshal(bytes)
	if err != nil {
		return nil, err
	}

	if len(jwks.Keys) == 0 {
		return nil, errors.New("auth server returned no keys")
	}

	keys := make(map[string]*rsa.PublicKey, len(jwks.Keys))
	for _, key := range jwks.Keys {
		publicKey, err := key.DecodePublicKey()
		if err != nil {
			return nil, err
		}

		keys[key.Kid] = publicKey.(*rsa.PublicKey)
	}

	return keys, nil
}
