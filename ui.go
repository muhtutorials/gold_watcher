package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"time"
)

func (app *Config) makeUI() {
	openPrice, currentPrice, priceChange := app.GetPriceText()

	priceContent := container.NewGridWithColumns(3, openPrice, currentPrice, priceChange)

	app.PriceContainer = priceContent

	toolbar := app.getToolbar()
	app.Toolbar = toolbar

	priceTabContent := app.pricesTab()
	holdingsTabContent := app.holdingsTab()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(
			"Prices",
			theme.HomeIcon(),
			priceTabContent,
		),
		container.NewTabItemWithIcon(
			"Holdings",
			theme.InfoIcon(),
			holdingsTabContent,
		),
	)
	tabs.SetTabLocation(container.TabLocationTop) // top is default

	finalContent := container.NewVBox(priceContent, toolbar, tabs)

	app.MainWindow.SetContent(finalContent)

	go func() {
		for range time.Tick(time.Second * 30) {
			app.refreshPriceContent()
		}
	}()
}

func (app *Config) refreshPriceContent() {
	openPrice, currentPrice, priceChange := app.GetPriceText()
	app.PriceContainer.Objects = []fyne.CanvasObject{openPrice, currentPrice, priceChange}
	app.PriceContainer.Refresh()

	chart := app.getChart()
	app.PriceChartContainer.Objects = []fyne.CanvasObject{chart}
	app.PriceChartContainer.Refresh()
}

func (app *Config) refreshHoldingsTable() {
	app.Holdings = app.getHoldingSlice()
	app.HoldingsTable.Refresh()
}
