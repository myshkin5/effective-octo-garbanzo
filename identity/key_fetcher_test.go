package identity_test

//go:generate hel

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/myshkin5/effective-octo-garbanzo/identity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyFetcher", func() {
	var (
		mockHTTPClient *mockHTTPClient
	)

	BeforeEach(func() {
		mockHTTPClient = newMockHTTPClient()
	})

	Describe("MustFetchKeys", func() {
		It("eventually succeeds", func() {
			mockHTTPClient.GetOutput.Resp <- nil
			mockHTTPClient.GetOutput.Err <- errors.New("bad things happened")

			mockHTTPClient.GetOutput.Resp <- createValidResponse("RSA", validModulus)
			mockHTTPClient.GetOutput.Err <- nil

			keys := identity.MustFetchKeys("http://somewhere.com", mockHTTPClient)

			Expect(keys["key1"]).NotTo(BeNil())
			Expect(keys["key1"].E).To(Equal(65537))
			Expect(keys["key2"]).NotTo(BeNil())
			Expect(keys["key2"].E).To(Equal(65538))
		})
	})

	Describe("FetchKeys", func() {
		Context("error response from client", func() {
			BeforeEach(func() {
				mockHTTPClient.GetOutput.Resp <- nil
				mockHTTPClient.GetOutput.Err <- errors.New("bad things happened")
			})

			It("attempts to get the resource on invocation", func() {
				identity.FetchKeys("http://somewhere.com", mockHTTPClient)

				Expect(mockHTTPClient.GetCalled).To(Receive(Equal(true)))
				Expect(mockHTTPClient.GetInput.Url).To(Receive(Equal("http://somewhere.com")))
			})

			It("returns an error when the client returns an error", func() {
				_, err := identity.FetchKeys("http://somewhere.com", mockHTTPClient)

				Expect(err).To(MatchError("bad things happened"))
			})
		})

		It("returns an error when the client returns a non-200 response code", func() {
			mockHTTPClient.GetOutput.Resp <- &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(nil),
			}
			mockHTTPClient.GetOutput.Err <- nil

			_, err := identity.FetchKeys("http://somewhere.com", mockHTTPClient)

			Expect(err).To(MatchError("auth server returned a non-200 response code, code was 500"))
		})

		It("returns an error if the body of the response doesn't parse", func() {
			mockHTTPClient.GetOutput.Resp <- &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`not-json`)),
			}
			mockHTTPClient.GetOutput.Err <- nil

			_, err := identity.FetchKeys("http://somewhere.com", mockHTTPClient)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid character"))
		})

		It("returns an error if there are no keys", func() {
			mockHTTPClient.GetOutput.Resp <- &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`{
						"keys": []
					}`)),
			}
			mockHTTPClient.GetOutput.Err <- nil

			_, err := identity.FetchKeys("http://somewhere.com", mockHTTPClient)

			Expect(err).To(MatchError("auth server returned no keys"))
		})

		It("returns an error if the key type is bogus", func() {
			mockHTTPClient.GetOutput.Resp <- createValidResponse("BogusKty", validModulus)
			mockHTTPClient.GetOutput.Err <- nil

			_, err := identity.FetchKeys("http://somewhere.com", mockHTTPClient)

			Expect(err).To(MatchError("Unknown JWK key type BogusKty"))
		})

		It("returns an error if there is a malformed modulus", func() {
			mockHTTPClient.GetOutput.Resp <- createValidResponse("RSA", "bogus-modulus")
			mockHTTPClient.GetOutput.Err <- nil

			_, err := identity.FetchKeys("http://somewhere.com", mockHTTPClient)

			Expect(err).To(MatchError("Malformed JWK RSA key"))
		})

		It("returns the public key", func() {
			mockHTTPClient.GetOutput.Resp <- createValidResponse("RSA", validModulus)
			mockHTTPClient.GetOutput.Err <- nil

			keys, err := identity.FetchKeys("http://somewhere.com", mockHTTPClient)

			Expect(err).NotTo(HaveOccurred())
			Expect(keys["key1"]).NotTo(BeNil())
			Expect(keys["key1"].E).To(Equal(65537))
			Expect(keys["key2"]).NotTo(BeNil())
			Expect(keys["key2"].E).To(Equal(65538))
		})
	})
})

const validModulus = "r0EA8mqy0JG7SpjTzPqt2Arp50CcxxVFbirmWRx1ELY2zFVc9g7O8j" +
	"L8ubBnJogojUVhl_hqBoiNDDuLmwyIfPbD1KZJ8vWhci0HMyJSFoHj8cayT_vliYlPXqmH_" +
	"pGh_WLI33VJ7mt_q5Ou_AltDXIPIkUdnKcSGKzDGSmIgroul3g8N8yPUIrCHp1vGiDqUguE" +
	"odUFk3h94RiBtdGXIEGyA9w1P2FpfyIzPPnGcRIknv_-DD7chqeUzBnAofXKsBfsvtWT_uX" +
	"vO4-umz43i0iZsRW4BRUitEZqxLzhL5h6JsLDCCpBWQzSDVg-_eKmq9_lJKdZJc0stUv-Xd" +
	"EkHQ"

func createValidResponse(keyType, modulus string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(strings.NewReader(`{
					"keys": [
						{
							"kid": "key1",
							"kty": "` + keyType + `",
							"n": "` + modulus + `",
							"e": "AQAB"
						},
						{
							"kid": "key2",
							"kty": "` + keyType + `",
							"n": "` + modulus + `",
							"e": "AQAC"
						}
					]
				}`)),
	}
}
