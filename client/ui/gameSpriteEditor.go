package ui

import (
	"fmt"
	"image"
	"io"
	"leafcli/models"
	"leafcli/ui/widgets"
	"leafcli/utils"
	"os"
	"path/filepath"

	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func buildGameSpriteEditor(p models.Project, w fyne.Window, spr models.Sprite) fyne.CanvasObject {
	name := widget.NewEntry()
	name.SetText(spr.Name)

	columns := widgets.NewNumberEntry()
	columns.MinValue = 1

	play_stop := widget.NewCheck("Play/Stop", nil)

	speed := widgets.NewNumberEntry()
	speed.MinValue = 0
	speed.ChangeValue(int(spr.AnimSpeed)) // default value

	var srcImg image.Image
	var imgcanvas *canvas.Image

	fmt.Println("spr.Columns:", spr.Columns)
	var imgSliced []*image.RGBA
	var frameIndex int
	var elapsed float64
	if spr.ImagePath != "" {
		file, _ := os.Open(spr.ImagePath)

		defer file.Close()

		// Décoder
		srcImg, _, _ = image.Decode(file)

		if spr.Columns != 0 {

			columns.ChangeValue(int(spr.Columns))
			imgSliced = utils.SplitImage(srcImg, 1, int(spr.Columns))
			utils.ChangeImage(&imgcanvas, imgSliced, 0)

		} else {
			imgcanvas = canvas.NewImageFromFile(spr.ImagePath)
			imgcanvas.SetMinSize(fyne.NewSize(float32(imgSliced[0].Bounds().Dx()), 64))
			imgcanvas.FillMode = canvas.ImageFillContain
			imgcanvas.ScaleMode = canvas.ImageScalePixels
			imgcanvas.Refresh()
		}

	} else {
		imgcanvas = &canvas.Image{}
	}

	imageSelect := widget.NewButton("Select Image", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, w)
				// navigation controls
				return
			}

			// Ensure sprite directory
			targetDir := spr.Directory
			if targetDir == "" {
				targetDir = filepath.Join(p.Directory, "Sprites")
			}
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				dialog.ShowError(err, w)
				return
			}

			destPath := filepath.Join(targetDir, filepath.Base(reader.URI().Path()))
			if err := os.WriteFile(destPath, data, 0644); err != nil {
				dialog.ShowError(err, w)
				return
			}

			spr.ImagePath = destPath

			file, err := os.Open(destPath)
			if err != nil {
				return
			}
			defer file.Close()

			// Décoder
			srcImg, _, err = image.Decode(file)
			if err != nil {
				return
			}
			imgSliced = utils.SplitImage(srcImg, 1, int(spr.Columns))

			utils.ChangeImage(&imgcanvas, imgSliced, 0)

			// Mettre à jour le conteneur avec la nouvelle image

			spr.ImageWidth = uint(srcImg.Bounds().Max.X)

			spr.ImageHeight = uint(srcImg.Bounds().Max.Y)

			if err := spr.Save(); err != nil {
				dialog.ShowError(err, w)
			}

		}, w).Show()
	})

	img := container.NewHBox(
		imageSelect,
		imgcanvas,
	)

	// boucle d'animation (goroutine) — utilise dt pour gérer la vitesse
	go func() {
		last := time.Now()
		ticker := time.NewTicker(time.Millisecond * 16)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			dt := now.Sub(last).Seconds()
			last = now
			if len(imgSliced) <= 1 || !play_stop.Checked {
				continue
			}
			fps := float64(speed.Value)
			if fps > 0 {

				frameDuration := 1.0 / fps
				elapsed += dt
				if elapsed >= frameDuration {
					elapsed -= frameDuration
					frameIndex = (frameIndex + 1) % len(imgSliced)
					utils.ChangeImage(&imgcanvas, imgSliced, frameIndex)

					img.Objects = []fyne.CanvasObject{imageSelect, imgcanvas}
					fyne.Do(img.Refresh)

				}
			}
		}
	}()

	anim := container.NewHBox(
		container.NewHBox(
			widget.NewLabel("Columns"),
			columns,
		),
		container.NewHBox(
			widget.NewLabel("Speed (FPS)"),
			speed,
		),
		container.NewHBox(play_stop),
	)

	saveBtn := widget.NewButton("Save", func() {
		fmt.Println("spr.Columns: ", spr.Columns)

		imgSliced = utils.SplitImage(srcImg, 1, int(columns.Value))
		utils.ChangeImage(&imgcanvas, imgSliced, 0)

		// Mettre à jour le conteneur avec la nouvelle image
		img.Objects = []fyne.CanvasObject{imageSelect, imgcanvas}
		img.Refresh()

		spr.AnimSpeed = uint(speed.Value)
		spr.Name = name.Text
		spr.Columns = uint(columns.Value)

		fmt.Println("test: spr: ", spr)
		if err := spr.Save(); err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Sprite", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Name", name),
			widget.NewFormItem("Animation", anim),
			widget.NewFormItem("Image", img),
		),
		saveBtn,
	)
}
