package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gold_watcher/repository"
	"strconv"
	"time"
)

func (app *Config) getToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addHoldingsDialog()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.refreshPriceContent()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			w := app.showPreferences()
			w.Resize(fyne.NewSize(300, 20))
			w.Show()
		}),
	)
	return toolbar
}

func (app *Config) addHoldingsDialog() dialog.Dialog {
	addAmountEntry := widget.NewEntry()
	purchaseDateEntry := widget.NewEntry()
	purchasePriceEntry := widget.NewEntry()

	app.AddHoldingsPurchaseAmountEntry = addAmountEntry
	app.AddHoldingsPurchaseDateEntry = purchaseDateEntry
	app.AddHoldingsPurchasePriceEntry = purchasePriceEntry

	dateValidator := func(s string) error {
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err
		}
		return nil
	}
	purchaseDateEntry.Validator = dateValidator

	isIntValidator := func(s string) error {
		if _, err := strconv.Atoi(s); err != nil {
			return err
		}
		return nil
	}
	addAmountEntry.Validator = isIntValidator

	isFloatValidator := func(s string) error {
		if _, err := strconv.ParseFloat(s, 32); err != nil {
			return err
		}
		return nil
	}
	purchasePriceEntry.Validator = isFloatValidator

	purchaseDateEntry.PlaceHolder = "YYYY-MM-DD"

	addForm := dialog.NewForm(
		"Add Gold",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Amount in toz", Widget: addAmountEntry},
			{Text: "Purchase Date", Widget: purchaseDateEntry},
			{Text: "Purchase Price", Widget: purchasePriceEntry},
		},
		func(valid bool) {
			if valid {
				amount, _ := strconv.Atoi(addAmountEntry.Text)
				purchaseDate, _ := time.Parse("2006-01-02", purchaseDateEntry.Text)
				purchasePrice, _ := strconv.ParseFloat(purchasePriceEntry.Text, 32)
				purchasePrice *= 100

				_, err := app.DB.InsertHolding(repository.Holdings{
					Amount:        amount,
					PurchaseDate:  purchaseDate,
					PurchasePrice: int(purchasePrice),
				})
				if err != nil {
					app.ErrorLog.Println(err)
				}
				app.refreshHoldingsTable()
			}
		},
		app.MainWindow,
	)

	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm
}

func (app *Config) showPreferences() fyne.Window {
	win := app.App.NewWindow("Preferences")

	lbl := widget.NewLabel("Preferred Currency")
	cur := widget.NewSelect([]string{"EUR", "RUB"}, func(s string) {
		currency = s
		app.App.Preferences().SetString("currency", s)
	})
	cur.Selected = currency

	btn := widget.NewButton("Save", func() {
		win.Close()
		app.refreshPriceContent()
	})
	btn.Importance = widget.HighImportance

	win.SetContent(container.NewVBox(lbl, cur, btn))

	return win
}
