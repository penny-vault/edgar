package financials_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/penny-vault/edgar/financials"
)

var _ = Describe("Statement", func() {

	var statement *financials.Statement

	BeforeEach(func() {
		var err error
		statement, err = financials.ParseXBRL("../testdata/jef.xml")
		Expect(err).To(BeNil())
	})

	Describe("CompanyProfile", func() {

		Context("form type", func() {
			It("should be 10-Q", func() {
				Expect(statement.FormType).To(Equal("10-Q"))
			})
		})

		Context("name", func() {
			It("should be Jefferies", func() {
				Expect(statement.CompanyProfile.Name).To(Equal("JEFFERIES FINANCIAL GROUP INC."))
			})
		})

		Context("ticker", func() {
			It("should be JEF", func() {
				Expect(statement.CompanyProfile.PrimaryTicker).To(Equal("JEF"))
			})
		})

		Context("exchange", func() {
			It("should be NYSE", func() {
				Expect(statement.CompanyProfile.PrimaryExchange).To(Equal("NYSE"))
			})
		})

		Context("total assets", func() {
			It("should be 52 billion", func() {
				Expect(statement.BalanceSheet.TotalAssets).To(BeNumerically("~", 5.20329e+10))
			})
		})

		Context("total liabilities", func() {
			It("should be 494 million", func() {
				Expect(statement.BalanceSheet.TotalLiabilities).To(BeNumerically("~", 4.945e+08))
			})
		})

		Context("current assets", func() {
			It("should be 30 billion", func() {
				Expect(statement.BalanceSheet.CurrentAssets).To(BeNumerically("~", 30_027_440_000))
			})
		})

		Context("current liabilities", func() {
			It("should be 9.44 billion", func() {
				Expect(statement.BalanceSheet.CurrentAssets).To(BeNumerically("~", 9_448_900_000))
			})
		})

		Context("cash and equivalents", func() {
			It("should be 7.44 billion", func() {
				Expect(statement.BalanceSheet.CashAndEquiv).To(BeNumerically("~", 7_508_508_000))
			})
		})

	})

})
