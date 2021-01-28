package exchangerates

import (
	"encoding/xml"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
)

const euroCentralBankExchangeRateUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
const euroCentralBankExchangeRateReferenceUrl = "https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html"
const euroCentralBankDataSource = "European Central Bank"
const euroCentralBankBaseCurrency = "EUR"

// EuroCentralBankDataSource defines the structure of exchange rates data source of euro central bank
type EuroCentralBankDataSource struct {
	ExchangeRatesDataSource
}

// EuroCentralBankExchangeRateData represents the whole data from euro central bank
type EuroCentralBankExchangeRateData struct {
	XMLName          xml.Name                        `xml:"Envelope"`
	AllExchangeRates []*EuroCentralBankExchangeRates `xml:"Cube>Cube"`
}

// EuroCentralBankExchangeRates represents the exchange rates data from euro central bank
type EuroCentralBankExchangeRates struct {
	Date          string                         `xml:"time,attr"`
	ExchangeRates []*EuroCentralBankExchangeRate `xml:"Cube"`
}

// EuroCentralBankExchangeRate represents the exchange rate data from euro central bank
type EuroCentralBankExchangeRate struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from euro central bank
func (e *EuroCentralBankExchangeRateData) ToLatestExchangeRateResponse() *models.LatestExchangeRateResponse {
	if len(e.AllExchangeRates) < 1 {
		return nil
	}

	latestEuroCentralBankExchangeRate := e.AllExchangeRates[0]

	if len(latestEuroCentralBankExchangeRate.ExchangeRates) < 1 {
		return nil
	}

	exchangeRates := make([]*models.LatestExchangeRate, len(latestEuroCentralBankExchangeRate.ExchangeRates))

	for i := 0; i < len(latestEuroCentralBankExchangeRate.ExchangeRates); i++ {
		exchangeRates[i] = latestEuroCentralBankExchangeRate.ExchangeRates[i].ToLatestExchangeRate()
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    euroCentralBankDataSource,
		ReferenceUrl:  euroCentralBankExchangeRateReferenceUrl,
		Date:          latestEuroCentralBankExchangeRate.Date,
		BaseCurrency:  euroCentralBankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from euro central bank
func (e *EuroCentralBankExchangeRate) ToLatestExchangeRate() *models.LatestExchangeRate {
	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     e.Rate,
	}
}

// GetRequestUrl returns the euro central bank data source url
func (e *EuroCentralBankDataSource) GetRequestUrl() string {
	return euroCentralBankExchangeRateUrl
}

// Parse returns the common response entity according to the euro central bank data source raw response
func (e *EuroCentralBankDataSource) Parse(c *core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	euroCentralBankData := &EuroCentralBankExchangeRateData{}
	err := xml.Unmarshal(content, euroCentralBankData)

	if err != nil {
		log.ErrorfWithRequestId(c, "[euro_central_bank_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := euroCentralBankData.ToLatestExchangeRateResponse()

	if latestExchangeRateResponse == nil {
		log.ErrorfWithRequestId(c, "[euro_central_bank_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse.ExchangeRates = append(latestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: euroCentralBankBaseCurrency,
		Rate:     "1",
	})

	return latestExchangeRateResponse, nil
}