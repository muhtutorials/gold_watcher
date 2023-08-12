package main

import (
	"fyne.io/fyne/v2/test"
	"testing"
)

func Test_GetToolbar(t *testing.T) {
	tb := testApp.getToolbar()

	if len(tb.Items) != 4 {
		t.Error("Wrong number of items in toolbar")
	}
}

func Test_AddHoldingsDialog(t *testing.T) {
	testApp.addHoldingsDialog()
	test.Type(testApp.AddHoldingsPurchaseAmountEntry, "1")
	test.Type(testApp.AddHoldingsPurchaseDateEntry, "1000")
	test.Type(testApp.AddHoldingsPurchasePriceEntry, "2020-11-05")

	if testApp.AddHoldingsPurchaseAmountEntry.Text != "1" {
		t.Error("amount is not correct")
	}

	if testApp.AddHoldingsPurchaseDateEntry.Text != "2020-11-05" {
		t.Error("date is not correct")
	}

	if testApp.AddHoldingsPurchasePriceEntry.Text != "1000" {
		t.Error("price is not correct")
	}
}
