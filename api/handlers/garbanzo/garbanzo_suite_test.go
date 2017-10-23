package garbanzo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGarbanzo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API - Handlers - Garbanzo Suite")
}
