package middleware_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API - Handlers - Middleware Suite")
}
