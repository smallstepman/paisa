package ledger

import (
	"testing"

	"github.com/ananthakumaran/paisa/internal/model/price"
	"github.com/ananthakumaran/paisa/internal/utils"
	"github.com/stretchr/testify/assert"
)

func assertPriceEqual(t *testing.T, actual price.Price, date string, commodityName string, value float64) {
	assert.Equal(t, commodityName, actual.CommodityName, "they should be equal")
	assert.Equal(t, date, actual.Date.Format("2006/01/02"), "they should be equal")
	assert.Equal(t, value, actual.Value.InexactFloat64(), "they should be equal")
}

func TestParseLegerPrices(t *testing.T) {
	parsedPrices, _ := parseLedgerPrices("P 2023/05/01 00:00:00 USD 0.9 EUR\n", "EUR")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 0.9)
	parsedPrices, _ = parseLedgerPrices("P 2023/05/01 00:00:00 EUR $1.1\n", "$")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "EUR", 1.1)
	parsedPrices, _ = parseLedgerPrices("P 2023/05/01 00:00:00 EUR $-1.1\n", "$")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "EUR", -1.1)
	parsedPrices, _ = parseLedgerPrices("P 2023/05/01 00:00:00 EUR ₹70\n", "₹")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "EUR", 70)

	parsedPrices, _ = parseLedgerPrices("P 2023/05/01 00:00:00 USD 0.9 EUR\n", "INR")
	assert.Len(t, parsedPrices, 0)
	parsedPrices, _ = parseLedgerPrices("P 2023/05/01 00:00:00 USD $0.9\n", "INR")
	assert.Len(t, parsedPrices, 0)

	parsedPrices, _ = parseLedgerPrices("P 2022/01/29 00:50:00 UAH 0.026 EUR\n", "EUR")
	assertPriceEqual(t, parsedPrices[0], "2022/01/29", "UAH", 0.026)
}

func TestParseHLegerPrices(t *testing.T) {
	parsedPrices, _ := parseHLedgerPrices("P 2023-05-01 USD 0.9 EUR\n", "EUR")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 0.9)
	parsedPrices, _ = parseHLedgerPrices("P 2023-05-01 EUR $1.1\n", "$")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "EUR", 1.1)

	parsedPrices, _ = parseHLedgerPrices("P 2023-05-01 EUR USD 1.1\n", "USD")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "EUR", 1.1)

	parsedPrices, _ = parseHLedgerPrices("P 2023-05-01 EUR 1.1$\n", "$")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "EUR", 1.1)

	parsedPrices, _ = parseHLedgerPrices(utils.Dos2Unix("P 2023-05-01 EUR 1.1$\r\n"), "$")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "EUR", 1.1)

	parsedPrices, _ = parseHLedgerPrices("P 2023-05-01 \"AAPL0\" \"USD0\" 45.5\n", "USD0")
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "AAPL0", 45.5)

	parsedPrices, _ = parseHLedgerPrices("P 2023-05-01 USD 0.9 EUR\n", "INR")
	assert.Len(t, parsedPrices, 0)

	parsedPrices, _ = parseHLedgerPrices("P 2023-05-01 USD $0.9\n", "INR")
	assert.Len(t, parsedPrices, 0)

	parsedPrices, _ = parseHLedgerPrices("P 2023-05-01 USD $0.9\r\n", "INR")
	assert.Len(t, parsedPrices, 0)
}

func TestParseAmount(t *testing.T) {
	commodity, amount, _ := parseAmount("0.9 USD")
	assert.Equal(t, "USD", commodity)
	assert.Equal(t, 0.9, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("$0.9")
	assert.Equal(t, "$", commodity)
	assert.Equal(t, 0.9, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("0.9$")
	assert.Equal(t, "$", commodity)
	assert.Equal(t, 0.9, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("$-0.9")
	assert.Equal(t, "$", commodity)
	assert.Equal(t, -0.9, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("-0.9$")
	assert.Equal(t, "$", commodity)
	assert.Equal(t, -0.9, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("100,000 EUR")
	assert.Equal(t, "EUR", commodity)
	assert.Equal(t, 100000.0, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("100,000.00 \"EUR0-0\"")
	assert.Equal(t, "EUR0-0", commodity)
	assert.Equal(t, 100000.0, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("-100,000.00 \"EUR0-0\"")
	assert.Equal(t, "EUR0-0", commodity)
	assert.Equal(t, -100000.0, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("\"EUR0-0\" -100,000.00")
	assert.Equal(t, "EUR0-0", commodity)
	assert.Equal(t, -100000.0, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("INR 70.0099")
	assert.Equal(t, "INR", commodity)
	assert.Equal(t, 70.0099, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("1E-8 BTC")
	assert.Equal(t, "BTC", commodity)
	assert.Equal(t, 1e-08, amount.InexactFloat64())

	commodity, amount, _ = parseAmount("100E-8    BTC")
	assert.Equal(t, "BTC", commodity)
	assert.Equal(t, 1e-06, amount.InexactFloat64())
}

func TestParseBeancountPrices(t *testing.T) {
	// Test CSV format with header
	csvWithHeader := `date,currency,amount_number,amount_currency
2023-05-01,USD,0.9,EUR
2023-05-02,EUR,1.1,USD
`
	parsedPrices, _ := parseBeancountPrices(csvWithHeader, "EUR")
	assert.Len(t, parsedPrices, 2)
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 0.9)
	assertPriceEqual(t, parsedPrices[1], "2023/05/02", "USD", 0.9090909090909091)

	// Test CSV format without header
	csvWithoutHeader := `2023-05-01,USD,0.9,EUR
2023-05-02,NIFTY,100.273,INR
`
	parsedPrices, _ = parseBeancountPrices(csvWithoutHeader, "EUR")
	assert.Len(t, parsedPrices, 1)
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 0.9)

	// Test normal case (target currency matches default)
	csvNormal := `2023-05-01,USD,80.442048,INR
`
	parsedPrices, _ = parseBeancountPrices(csvNormal, "INR")
	assert.Len(t, parsedPrices, 1)
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 80.442048)

	// Test filtering out non-matching currencies
	csvNonMatching := `2023-05-01,USD,0.9,EUR
2023-05-02,EUR,1.1,GBP
`
	parsedPrices, _ = parseBeancountPrices(csvNonMatching, "INR")
	assert.Len(t, parsedPrices, 0)

	// Test with quoted values
	csvQuoted := `2023-05-01,"USD",0.9,"EUR"
`
	parsedPrices, _ = parseBeancountPrices(csvQuoted, "EUR")
	assert.Len(t, parsedPrices, 1)
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 0.9)

	// Test with different header case variations
	csvMixedCase := `Date,Currency,Amount_Number,Amount_Currency
2023-05-01,USD,0.9,EUR
`
	parsedPrices, _ = parseBeancountPrices(csvMixedCase, "EUR")
	assert.Len(t, parsedPrices, 1)
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 0.9)

	csvUpperCase := `DATE,CURRENCY,AMOUNT_NUMBER,AMOUNT_CURRENCY
2023-05-01,USD,0.9,EUR
`
	parsedPrices, _ = parseBeancountPrices(csvUpperCase, "EUR")
	assert.Len(t, parsedPrices, 1)
	assertPriceEqual(t, parsedPrices[0], "2023/05/01", "USD", 0.9)
}

func TestBeancountPrices_FixtureScenario(t *testing.T) {
	// This is what bean-query returns for the inr-beancount fixture
	csvOutput := `date,currency,amount_number,amount_currency
2022-01-07,NIFTY,100,INR
2022-02-07,NIFTY,100.273,INR
2022-01-08,USD,80.442048,INR`

	parsedPrices, err := parseBeancountPrices(csvOutput, "INR")
	assert.NoError(t, err)
	assert.Len(t, parsedPrices, 3, "Should parse exactly 3 prices")

	// Count by commodity
	niftyPrices := 0
	usdPrices := 0
	for _, p := range parsedPrices {
		if p.CommodityName == "NIFTY" {
			niftyPrices++
		}
		if p.CommodityName == "USD" {
			usdPrices++
		}
	}

	assert.Equal(t, 2, niftyPrices, "Should have 2 NIFTY prices")
	assert.Equal(t, 1, usdPrices, "Should have 1 USD price")

	// Verify specific values match expected price.json
	assert.Equal(t, "NIFTY", parsedPrices[0].CommodityName)
	assert.Equal(t, 100.0, parsedPrices[0].Value.InexactFloat64())
	assert.Equal(t, "2022-01-07", parsedPrices[0].Date.Format("2006-01-02"))

	assert.Equal(t, "NIFTY", parsedPrices[1].CommodityName)
	assert.Equal(t, 100.273, parsedPrices[1].Value.InexactFloat64())
	assert.Equal(t, "2022-02-07", parsedPrices[1].Date.Format("2006-01-02"))

	assert.Equal(t, "USD", parsedPrices[2].CommodityName)
	assert.Equal(t, 80.442048, parsedPrices[2].Value.InexactFloat64())
	assert.Equal(t, "2022-01-08", parsedPrices[2].Date.Format("2006-01-02"))
}
