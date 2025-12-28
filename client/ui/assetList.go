package ui

import (
	"fmt"
	"leafcli/models"
	"leafcli/utils"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createCategory(edStatus *EditorStatus, categories []*widget.List, p models.Project, w fyne.Window, name string) fyne.CanvasObject {

	var categorie_created int = -1
	switch name {
	case "GameObjects":
		categorie_created = 0
		categories[0] = widget.NewList(
			func() int { return len(p.GameObjects) },
			func() fyne.CanvasObject { return widget.NewLabel("Asset") },
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(models.LoadGameObjects(p.GameObjects)[i].Name)
			},
		)

	case "Sprites":
		categorie_created = 1
		categories[1] = widget.NewList(
			func() int { return len(p.Sprites) },
			func() fyne.CanvasObject { return widget.NewLabel("Asset") },
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(models.LoadSprites(p.Sprites)[i].Name)
			},
		)

	case "Rooms":
		categorie_created = 2
		categories[2] = widget.NewList(
			func() int { return len(p.Rooms) },
			func() fyne.CanvasObject { return widget.NewLabel("Asset") },
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(models.LoadRooms(p.Rooms)[i].Name)

			},
		)
	}
	categories[categorie_created].OnSelected = func(id widget.ListItemID) {
		fmt.Println("Selected", name, "ID:", id)

		edStatus.selected.cat = name
		edStatus.selected.id = id

		// Clear editor
		edStatus.editorPanel.RemoveAll()
		edStatus.editorPanel.Refresh()

		switch name {
		case "GameObjects":
			gO := models.LoadGameObjects(p.GameObjects)[id]
			edStatus.editorPanel.Add(buildGameObjectEditor(p, w, gO))
			categories[1].UnselectAll()
			categories[2].UnselectAll()

		case "Sprites":
			sp := models.LoadSprites(p.Sprites)[id]
			edStatus.editorPanel.Add(buildGameSpriteEditor(p, w, sp))
			categories[0].UnselectAll()
			categories[2].UnselectAll()

		case "Rooms":
			roo := models.LoadRooms(p.Rooms)[id]
			edStatus.editorPanel.Add(buildRoomEditor(p, w, roo))
			categories[0].UnselectAll()
			categories[1].UnselectAll()
		}

		edStatus.editorPanel.Refresh()
		fmt.Println("edStatus.selected:", edStatus.selected)

	}

	scroll := container.NewVScroll(categories[categorie_created])
	scroll.SetMinSize(fyne.NewSize(200, 120)) // hauteur par défaut

	scroll.Hide() // Liste cachée au départ

	// Header = Label à gauche, boutons à droite
	addBtn := widget.NewButton("+", func() {
		ShowCreateAssetDialog(w, name, func(assetName string) {
			switch name {
			case "Sprites":
				var sp models.Sprite
				sp.Name = assetName
				sp.Directory = filepath.Join(p.Directory, "Sprites", sp.Name)
				sp.Columns = 1
				utils.EnsureDir(sp.Directory)
				p.Sprites = append(p.Sprites, filepath.Join(sp.Directory, "info.lsp"))

				sp.Save()
				fmt.Println("Create sprite:", assetName)

			case "GameObjects":
				var gO models.GameObject
				gO.Name = assetName

				gO.Directory = filepath.Join(p.Directory, "GameObjects", gO.Name)
				utils.EnsureDir(gO.Directory)
				p.GameObjects = append(p.GameObjects, filepath.Join(gO.Directory, "info.lgo"))
				gO.Scripts = make(map[string]string)
				gO.Save()

				fmt.Println("Create game object:", assetName)

			case "Rooms":
				var roo models.Room
				roo.Name = assetName
				roo.Directory = filepath.Join(p.Directory, "Rooms")
				utils.EnsureDir(roo.Directory)
				p.Rooms = append(p.Rooms, filepath.Join(roo.Directory, roo.Name+".lro"))
				roo.Save()
				fmt.Println("Create room:", assetName)
			}

			p.UpdateProject()
		})
	})
	toggleBtn := widget.NewButton("▾", func() {
		if name == edStatus.selected.cat {
			edStatus.selected.cat = ""
			edStatus.selected.id = 0
		}
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
