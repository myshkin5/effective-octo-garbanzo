package identity_test

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/myshkin5/effective-octo-garbanzo/identity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	var (
		privateKey *rsa.PrivateKey
		publicKeys map[string]*rsa.PublicKey
		validator  *identity.Validator
	)

	BeforeSuite(func() {
		var err error
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		Expect(err).NotTo(HaveOccurred())

		publicKeys = make(map[string]*rsa.PublicKey)
		publicKeys["joe"] = &privateKey.PublicKey
	})

	BeforeEach(func() {
		validator = identity.NewValidator(publicKeys)
	})

	createJWT := func(keyId interface{}, validSeconds int64) string {
		claims := jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + validSeconds,
		}
		token := &jwt.Token{
			Header: map[string]interface{}{
				"kid": keyId,
				"alg": jwt.SigningMethodRS256.Alg(),
			},
			Claims: claims,
			Method: jwt.SigningMethodRS256,
		}
		signedString, err := token.SignedString(privateKey)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		return signedString
	}

	It("reports bogus headers as invalid", func() {
		ok := validator.IsValid("bogus")
		Expect(ok).To(BeFalse())
	})

	It("reports bogus tokens as invalid", func() {
		ok := validator.IsValid("bearer bogus")
		Expect(ok).To(BeFalse())
	})

	It("reports good tokens as valid regardless of the prefix case", func() {
		ok := validator.IsValid("Bearer " + createJWT("joe", 5))
		Expect(ok).To(BeTrue())
	})

	It("reports non-string key ids as invalid", func() {
		ok := validator.IsValid("Bearer " + createJWT(22, 5))
		Expect(ok).To(BeFalse())
	})

	It("reports expired tokens as invalid", func() {
		ok := validator.IsValid("Bearer " + createJWT("joe", -5))
		Expect(ok).To(BeFalse())
	})

	It("reports good tokens with no matching key as invalid", func() {
		ok := validator.IsValid("Bearer " + createJWT("alice", 5))
		Expect(ok).To(BeFalse())
	})
})
