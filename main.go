package main

import (
	"database/sql"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	_ "github.com/glebarez/go-sqlite"
	"gold_watcher/repository"
	"log"
	"net/http"
	"os"
)

type Config struct {
	App                            fyne.App
	InfoLog                        *log.Logger
	ErrorLog                       *log.Logger
	DB                             repository.Repository
	MainWindow                     fyne.Window
	PriceContainer                 *fyne.Container
	Toolbar                        *widget.Toolbar
	PriceChartContainer            *fyne.Container
	Holdings                       [][]any
	HoldingsTable                  *widget.Table
	HTTPClient                     *http.Client
	AddHoldingsPurchaseAmountEntry *widget.Entry
	AddHoldingsPurchaseDateEntry   *widget.Entry
	AddHoldingsPurchasePriceEntry  *widget.Entry
}

func main() {
	var myApp Config

	fyneApp := app.NewWithID("ca.gocode.goldwatcher.preferences")
	myApp.App = fyneApp
	myApp.HTTPClient = &http.Client{}
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	myApp.setupDB(sqlDB)

	currency = fyneApp.Preferences().StringWithFallback("currency", "USD")

	myApp.MainWindow = fyneApp.NewWindow("Gold Watcher")
	myApp.MainWindow.Resize(fyne.NewSize(770, 410))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()

	myApp.makeUI()

	myApp.MainWindow.ShowAndRun()
}

func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Println(err)
		log.Panic()
	}
}
