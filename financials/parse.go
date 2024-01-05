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

import (
	"encoding/xml"
	"os"
	"strings"
	"time"

	"github.com/polygon-io/xbrl-parser"
	"github.com/rs/zerolog/log"
)

func ParseXBRL(fn string) (*Statement, error) {
	statement := &Statement{
		ProcessingDate: time.Now(),
	}

	var processed xbrl.XBRL
	doc, err := os.ReadFile(fn)
	if err != nil {
		log.Error().Err(err).Msg("read file failed")
		return nil, err
	}

	if err := xml.Unmarshal([]byte(doc), &processed); err != nil {
		log.Error().Err(err).Msg("unmarshal failed")
		return nil, err
	}

	// get end date
	endDate := ""
	for _, fact := range processed.Facts {
		if fact.IsValid() && fact.XMLName.Local == "DocumentPeriodEndDate" {
			endDate = *fact.ValueStr
		}
	}

	log.Info().Int("NumFacts", len(processed.Facts)).Msg("Parsed XBRL")
	for _, fact := range processed.Facts {
		if fact.IsValid() {
			switch fact.XMLName.Local {
			case "EntityRegistrantName":
				statement.CompanyProfile.Name = *fact.ValueStr
			case "SecurityExchangeName":
				statement.CompanyProfile.PrimaryExchange = *fact.ValueStr
			case "TradingSymbol":
				symbol := *fact.ValueStr
				if !strings.Contains(symbol, " ") && statement.CompanyProfile.PrimaryTicker == "" {
					statement.CompanyProfile.PrimaryTicker = symbol
				}
			case "DocumentType":
				statement.FormType = *fact.ValueStr
			case "Assets":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.TotalAssets, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing assets numeric value")
						return nil, err
					}
				}
			case "Liabilities":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.TotalLiabilities, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing liabilities numeric value")
						return nil, err
					}
				}
			case "CashAndCashEquivalentsAtCarryingValue":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.CashAndEquiv, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing current loans receivable numeric value")
						return nil, err
					}
				}
			case "InvestmentsInAffiliatesSubsidiariesAssociatesAndJointVentures":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.CurrentLoansReceivable, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing cash and equivalents numeric value")
						return nil, err
					}
				}
			case "ReceivablesFromBrokersDealersAndClearingOrganizations":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.ReceivablesFromBrokers, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing receivables from brokers dealers and clearing organizations numeric value")
						return nil, err
					}
				}
			case "ReceivablesFromCustomers":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.ReceivablesFromCustomers, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing receivables from customers numeric value")
						return nil, err
					}
				}
			case "FeesInterestAndOther":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.ReceivablesFromOther, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing fees interest and other numeric value")
						return nil, err
					}
				}
			case "FinancialInstrumentsOwnedAtFairValue":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.FinancialInstrumentsOwnedAtFairValue, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing financial instruments owned at fair value numeric value")
						return nil, err
					}
				}
			case "CashAndSecuritiesSegregatedUnderFederalAndOtherRegulations":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.SegregatedCash, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing financial instruments owned at fair value numeric value")
						return nil, err
					}
				}
			case "SecuritiesReceivedAsCollateral":
				factContext := processed.ContextsByID[fact.ContextRef]
				if *factContext.Period.Instant == endDate {
					statement.BalanceSheet.SecuritiesReceivedAsCollateral, err = fact.NumericValue()
					if err != nil {
						log.Error().Err(err).Msg("error parsing financial instruments owned at fair value numeric value")
						return nil, err
					}
				}
			}
		}
	}

	if statement.BalanceSheet.CurrentAssets == 0 {
		statement.BalanceSheet.CurrentAssets = (statement.BalanceSheet.CashAndEquiv +
			statement.BalanceSheet.CurrentLoansReceivable +
			statement.BalanceSheet.ReceivablesFromBrokers +
			statement.BalanceSheet.ReceivablesFromCustomers +
			statement.BalanceSheet.ReceivablesFromOther +
			statement.BalanceSheet.FinancialInstrumentsOwnedAtFairValue +
			statement.BalanceSheet.SegregatedCash +
			statement.BalanceSheet.SecuritiesReceivedAsCollateral)
	}

	if statement.BalanceSheet.CurrentLiabilities == 0 {
		statement.BalanceSheet.CurrentLiabilities = statement.BalanceSheet.TotalLiabilities
	}

	statement.BalanceSheet.WorkingCapital = statement.BalanceSheet.CurrentAssets - statement.BalanceSheet.CurrentLiabilities

	return statement, nil
}
