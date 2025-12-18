package ui

import (
	"fmt"
	"leafcli/models"
	"leafcli/utils"
	"path"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowCreateAssetDialog(
	w fyne.Window,
	assetType string,
	onCreate func(name string),
) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Asset name")

	dialog.NewForm(
		"New "+assetType,
		"Create",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
		},
		func(confirm bool) {
			if !confirm {
				return
			}
			if nameEntry.Text == "" {
				dialog.ShowError(
					fmt.Errorf("asset name cannot be empty"),
					w,
				)
				return
			}
			onCreate(nameEntry.Text)
		},
		w,
	).Show()
}

func CreateProjectEditor(a fyne.App, p models.Project) fyne.Window {
	w := a.NewWindow("LeafCLI: " + p.Name + " " + p.Version.String())
	w.Resize(fyne.NewSize(1000, 600))
	var assetList *widget.List

	createCategory := func(name string) fyne.CanvasObject {
		switch name {
		case "GameObjects":
			assetList = widget.NewList(
				func() int { return len(p.GameObjects) },
				func() fyne.CanvasObject { return widget.NewLabel("Asset") },
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(models.LoadGameObjects(p.GameObjects)[i].Name)
				},
			)

		case "Sprites":
			assetList = widget.NewList(
				func() int { return len(p.Sprites) },
				func() fyne.CanvasObject { return widget.NewLabel("Asset") },
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(models.LoadSprites(p.Sprites)[i].Name)
				},
			)
		case "Rooms":
			assetList = widget.NewList(
				func() int { return len(p.Rooms) },
				func() fyne.CanvasObject { return widget.NewLabel("Asset") },
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(models.LoadRooms(p.Rooms)[i].Name)
				},
			)
		}
		scroll := container.NewVScroll(assetList)
		scroll.SetMinSize(fyne.NewSize(200, 120)) // hauteur par défaut

		scroll.Hide() // Liste cachée au départ

		// Header = Label à gauche, boutons à droite
		addBtn := widget.NewButton("+", func() {
			ShowCreateAssetDialog(w, name, func(assetName string) {
				switch name {
				case "Sprites":
					var sp models.Sprite
					sp.Name = assetName
					sp.Directory = path.Join(p.Directory, "Sprites")
					utils.EnsureDir(sp.Directory)
					p.GameObjects = append(p.Sprites, filepath.Join(sp.Directory, sp.Name+".lsp"))

					sp.Save()
					fmt.Println("Create sprite:", assetName)

				case "GameObjects":
					var gO models.GameObject
					gO.Name = assetName

					gO.Directory = path.Join(p.Directory, "GameObjects")
					utils.EnsureDir(gO.Directory)

					p.GameObjects = append(p.GameObjects, filepath.Join(gO.Directory, gO.Name+".lgo"))

					gO.Save()

					fmt.Println("Create game object:", assetName)

				case "Rooms":
					var roo models.Room
					roo.Name = assetName
					roo.Directory = path.Join(p.Directory, "Rooms")
					utils.EnsureDir(roo.Directory)
					p.GameObjects = append(p.Rooms, filepath.Join(roo.Directory, roo.Name+".lro"))
					roo.Save()
					fmt.Println("Create room:", assetName)
				}

				p.UpdateProject()
			})
		})
		toggleBtn := widget.NewButton("▾", func() {
			if scroll.Visible() {
				scroll.Hide()
			} else {
				scroll.Show()
			}
		})

		header := container.NewHBox(
			widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(), // pousse les boutons à droite
			addBtn,
			toggleBtn,
		)

		return container.NewVBox(header, scroll)
	}

	// ---- Categories
	spritesCat := createCategory("Sprites")
	gameObjectsCat := createCategory("GameObjects")
	roomsCat := createCategory("Rooms")

	leftPanel := container.NewVBox(
		spritesCat,
		gameObjectsCat,
		roomsCat,
	)

	// ---- Layout principal
	content := container.NewHSplit(
		leftPanel,
		widget.NewLabel("Editor area"), // placeholder pour la zone de travail
	)
	content.SetOffset(0.25)

	w.SetContent(content)
	return w
}
