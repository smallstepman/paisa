// +build integration

package ledger

import (
"os"
"path/filepath"
"testing"

"github.com/ananthakumaran/paisa/internal/config"
"github.com/stretchr/testify/assert"
)

func TestBeancountIntegration(t *testing.T) {
// Set up config for INR currency
fixtureDir, _ := filepath.Abs("../../tests/fixture/inr-beancount")
configPath := filepath.Join(fixtureDir, "paisa.yaml")
os.Setenv("PAISA_CONFIG", configPath)
defer os.Unsetenv("PAISA_CONFIG")

// Load config
config.LoadConfigFile(configPath)

journalPath := filepath.Join(fixtureDir, "main.beancount")

bc := Beancount{}

// Test ValidateFile
t.Run("ValidateFile", func(t *testing.T) {
errors, output, err := bc.ValidateFile(journalPath)
assert.NoError(t, err, "ValidateFile should not return error")
assert.Empty(t, errors, "Should have no validation errors")
assert.Equal(t, "", output, "Should return empty string after successful validation")
})

// Test Prices
t.Run("Prices", func(t *testing.T) {
prices, err := bc.Prices(journalPath)
assert.NoError(t, err, "Prices should not return error")
assert.Len(t, prices, 3, "Should extract exactly 3 prices")

// Count by commodity
niftyCount := 0
usdCount := 0
for _, p := range prices {
t.Logf("Price: %s = %.6f INR on %s", p.CommodityName, p.Value.InexactFloat64(), p.Date.Format("2006-01-02"))
if p.CommodityName == "NIFTY" {
niftyCount++
}
if p.CommodityName == "USD" {
usdCount++
}
}

assert.Equal(t, 2, niftyCount, "Should have 2 NIFTY prices")
assert.Equal(t, 1, usdCount, "Should have 1 USD price")

// Verify specific values
for _, p := range prices {
switch p.CommodityName {
case "NIFTY":
val := p.Value.InexactFloat64()
assert.True(t, val == 100.0 || val == 100.273, "NIFTY price should be 100 or 100.273")
case "USD":
assert.InDelta(t, 80.442048, p.Value.InexactFloat64(), 0.000001, "USD price should be 80.442048")
}
}
})
}
