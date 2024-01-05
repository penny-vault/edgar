package financials_test

import (
	"testing"

	"github.com/rs/zerolog/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFinancials(t *testing.T) {
	log.Logger = log.Output(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Financials Suite")
}
