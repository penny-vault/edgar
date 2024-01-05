// Copyright 2024
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package financials

import "time"

type Statement struct {
	FilingDate     time.Time
	CalendarDate   time.Time
	ProcessingDate time.Time

	FormType string

	CompanyProfile CompanyProfile

	IncomeStatement IncomeStatement
	BalanceSheet    BalanceSheet
	CashFlow        CashFlow
}

type CompanyProfile struct {
	CIK             string
	CUSIPs          []string
	PrimaryTicker   string
	PrimaryExchange string
	FIGI            string
	CompositeFIGI   string
	ShareClassFIGI  string
	Name            string
}

type IncomeStatement struct {
	Revenue                         float64
	CostOfRevenue                   float64
	SellingAndGeneralAdminExp       float64
	ResearchAndDevelExp             float64
	OPEX                            float64
	InterestExp                     float64
	TaxExp                          float64
	NetIncomeDiscontinuedOperations float64
	ConsolidatedIncome              float64
	NetIncomeNCI                    float64
	PreferredDividend               float64
	NetIncomeCommonStock            float64

	EPS        float64
	EPSDiluted float64

	SharesWeightedAvg        float64
	SharesWeightedAvgDiluted float64
}

type BalanceSheet struct {
	CashAndEquiv          float64
	TotalInvestments      float64
	CurrentInvestments    float64
	NonCurrentInvestments float64

	TotalAssets      float64
	CurrentAssets    float64
	NonCurrentAssets float64

	TotalLiabilities      float64
	CurrentLiabilities    float64
	NonCurrentLiabilities float64

	TotalDebt      float64
	CurrentDebt    float64
	NonCurrentDebt float64

	WorkingCapital float64

	DeferredRevenue          float64
	TotalDeposits            float64
	NetPropertyPlantAndEquip float64

	InventoryShortTerm float64

	TaxAssets      float64
	TaxLiabilities float64

	CurrentLoansReceivable   float64
	ReceivablesFromBrokers   float64
	ReceivablesFromCustomers float64
	ReceivablesFromOther     float64
	TotalReceivables         float64
	TotalPayables            float64
	Intangibles              float64

	Equity                              float64
	RetainedEarnings                    float64
	AccumulatedOtherComprehensiveIncome float64
}

type CashFlow struct {
	CAPEX                             float64
	NetBusinessAcqDivestures          float64
	NetInvestmentAcqDivestures        float64
	FreeCashFlowPerShare              float64
	NetCashFlowFromFinancing          float64
	TotalIssuanceRepaymentDebt        float64
	TotalIssuanceRepaymentEquity      float64
	CommonDividend                    float64
	NetCashFlowFromInvestments        float64
	NetCashFlowFromOperations         float64
	EffectOfForeignExchangeRateOnCash float64
	NetCashFlow                       float64
	StockBasedCompensation            float64
	DepreciationAmortization          float64
}
