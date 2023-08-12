package main

import "testing"

func TestConfig_GetPriceText(t *testing.T) {
	all, err := testApp.currentHoldings()
	if err != nil {
		t.Error("failed to get current from database:", err)
	}
	if len(all) != 2 {
		t.Error("wrong number of rows returned")
	}
}

func TestConfig_GetHoldingSlice(t *testing.T) {
	slice := testApp.getHoldingSlice()

	if len(slice) != 3 {
		t.Error("wrong number of rows returned; expected 3 but got", len(slice))
	}
}
