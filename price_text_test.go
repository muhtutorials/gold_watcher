package main

import "testing"

func TestApp_GetPriceText(t *testing.T) {
	open, _, _ := testApp.GetPriceText()
	if open.Text != "Open: $1935.7900 USD" {
		t.Error("Wrong price returned", open.Text)
	}
}
