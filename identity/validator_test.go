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

	type jwtRequest struct {
		keyId        interface{}
		algorithm    string
		method       jwt.SigningMethod
		validSeconds int64
		org          string
		signingKey   interface{}
	}

	createJWT := func(request jwtRequest) string {
		claims := identity.CustomClaims{
			Org: request.org,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + request.validSeconds,
			},
		}
		token := &jwt.Token{
			Header: map[string]interface{}{
				"kid": request.keyId,
				"alg": request.algorithm,
			},
			Claims: &claims,
			Method: request.method,
		}
		signedString, err := token.SignedString(request.signingKey)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())

		return signedString
	}

	It("reports bogus headers as invalid", func() {
		ok, _ := validator.IsValid("bogus")
		Expect(ok).To(BeFalse())
	})

	It("reports bogus tokens as invalid", func() {
		ok, _ := validator.IsValid("bearer bogus")
		Expect(ok).To(BeFalse())
	})

	It("reports good tokens as valid regardless of the prefix case", func() {
		ok, org := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        "joe",
			algorithm:    "RS256",
			method:       jwt.SigningMethodRS256,
			validSeconds: 500,
			org:          "org1",
			signingKey:   privateKey,
		}))
		Expect(ok).To(BeTrue())
		Expect(org).To(Equal("org1"))
	})

	It("reports non-string key ids as invalid", func() {
		ok, _ := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        22,
			algorithm:    "RS256",
			method:       jwt.SigningMethodRS256,
			validSeconds: 5,
			org:          "org1",
			signingKey:   privateKey,
		}))
		Expect(ok).To(BeFalse())
	})

	It("reports expired tokens as invalid", func() {
		ok, _ := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        "joe",
			algorithm:    "RS256",
			method:       jwt.SigningMethodRS256,
			validSeconds: -5,
			org:          "org1",
			signingKey:   privateKey,
		}))
		Expect(ok).To(BeFalse())
	})

	It("reports good tokens with no matching key as invalid", func() {
		ok, _ := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        "alice",
			algorithm:    "RS256",
			method:       jwt.SigningMethodRS256,
			validSeconds: 5,
			org:          "org1",
			signingKey:   privateKey,
		}))
		Expect(ok).To(BeFalse())
	})

	It("reports tokens with bad algorithms as invalid", func() {
		ok, _ := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        "joe",
			algorithm:    "none",
			method:       jwt.SigningMethodRS256,
			validSeconds: 5,
			org:          "org1",
			signingKey:   privateKey,
		}))
		Expect(ok).To(BeFalse())
	})

	It("reports tokens with bad signing methods as invalid", func() {
		ok, _ := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        "joe",
			algorithm:    "RS256",
			method:       jwt.SigningMethodNone,
			validSeconds: 5,
			org:          "org1",
			signingKey:   jwt.UnsafeAllowNoneSignatureType,
		}))
		Expect(ok).To(BeFalse())
	})

	It("reports tokens with bad signing methods and bad algorithms as invalid", func() {
		ok, _ := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        "joe",
			algorithm:    "none",
			method:       jwt.SigningMethodNone,
			validSeconds: 5,
			org:          "org1",
			signingKey:   jwt.UnsafeAllowNoneSignatureType,
		}))
		Expect(ok).To(BeFalse())
	})

	It("reports good tokens with no org claim as invalid", func() {
		ok, _ := validator.IsValid("Bearer " + createJWT(jwtRequest{
			keyId:        "joe",
			algorithm:    "RS256",
			method:       jwt.SigningMethodRS256,
			validSeconds: 5,
			signingKey:   privateKey,
		}))
		Expect(ok).To(BeFalse())
	})
})
