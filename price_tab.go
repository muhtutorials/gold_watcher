package main

import (
	"bytes"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image"
	"image/png"
	"io"
	"os"
	"strings"
)

func (app *Config) pricesTab() *fyne.Container {
	chart := app.getChart()
	chartContainer := container.NewVBox(chart)
	app.PriceChartContainer = chartContainer

	return chartContainer
}

func (app *Config) getChart() *canvas.Image {
	apiURL := fmt.Sprintf("https://goldprice.org/charts/gold_3d_b_o_%s_x.png", strings.ToLower(currency))
	var img *canvas.Image

	err := app.downloadFile(apiURL, "gold.png")
	if err != nil {
		// bundles image
		// fyne bundle -o bundled.go unreachable.png
		img = canvas.NewImageFromResource(resourceUnreachablePng)
	} else {
		img = canvas.NewImageFromFile("gold.png")
	}

	img.SetMinSize(fyne.Size{
		Width:  770,
		Height: 410,
	})
	img.FillMode = canvas.ImageFillOriginal

	return img
}

func (app *Config) downloadFile(URL, filename string) error {
	res, err := app.HTTPClient.Get(URL)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("Received wrong status code while downloading image.")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	out, err := os.Create(fmt.Sprintf("./%s", filename))

	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
