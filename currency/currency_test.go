package currency

import "testing"

func TestValidCurrency(t *testing.T) {
	c, _ := NewCode("GBP")
	if c.String() != "GBP" {
		t.Errorf("does not construct with valid code")
	}
}

func TestInvalidCurrency(t *testing.T) {
	_, err := NewCode("GB")
	if err.Error() != "invalid currency code" {
		t.Errorf("does not throw error on invalid code")
	}
}
