package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gold_watcher/repository"
	"strconv"
)

func (app *Config) holdingsTab() *fyne.Container {
	app.Holdings = app.getHoldingSlice()
	app.HoldingsTable = app.getHoldingsTable()

	holdingsContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, app.HoldingsTable),
	)

	return holdingsContainer
}

func (app *Config) getHoldingsTable() *widget.Table {
	t := widget.NewTable(
		func() (int, int) {
			return len(app.Holdings), len(app.Holdings[0])
		},
		func() fyne.CanvasObject {
			contr := container.NewVBox(widget.NewLabel(""))
			return contr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == (len(app.Holdings[0])-1) && i.Row != 0 {
				// last cell (button)
				w := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Delete?", "", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(app.Holdings[i.Row][0].(string))
							err := app.DB.DeleteHolding(int64(id))
							if err != nil {
								app.ErrorLog.Println(err)
							}
						}
						app.refreshHoldingsTable()
					}, app.MainWindow)
				})
				w.Importance = widget.HighImportance
				o.(*fyne.Container).Objects = []fyne.CanvasObject{w}
			} else {
				// text info
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(app.Holdings[i.Row][i.Col].(string)),
				}
			}
		})

	colWidth := []float32{50, 200, 200, 200, 110}
	for i := 0; i < len(colWidth); i++ {
		t.SetColumnWidth(i, colWidth[i])
	}

	return t
}

func (app *Config) getHoldingSlice() [][]any {
	var slice [][]any
	holdings, err := app.currentHoldings()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	slice = append(slice, []interface{}{"ID", "Amount", "Price", "Date", "Delete?"})

	for _, x := range holdings {
		var currentRow []any
		currentRow = append(currentRow, strconv.FormatInt(x.ID, 10))
		currentRow = append(currentRow, fmt.Sprintf("%d toz", x.Amount))
		currentAmount := float32(x.PurchasePrice) / 100
		currentRow = append(currentRow, fmt.Sprintf("$%.2f", currentAmount))
		currentRow = append(currentRow, x.PurchaseDate.Format("2006-01-02"))
		currentRow = append(currentRow, widget.NewButton("Delete", func() {

		}))

		slice = append(slice, currentRow)
	}
	return slice
}

func (app *Config) currentHoldings() ([]repository.Holdings, error) {
	holdings, err := app.DB.AllHoldings()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	return holdings, nil
}
