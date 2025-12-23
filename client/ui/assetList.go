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

func createCategory(edStatus *EditorStatus, assetList *widget.List, p models.Project, w fyne.Window, name string) fyne.CanvasObject {

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
	assetList.OnSelected = func(id widget.ListItemID) {
		edStatus.selected.cat = name
		edStatus.selected.id = id
		edStatus.editorPanel.Objects = nil // clear

		switch edStatus.selected.cat {
		case "GameObjects":
			var gO models.GameObject = models.LoadGameObjects(p.GameObjects)[edStatus.selected.id]
			edStatus.editorPanel.Add(buildGameObjectEditor(p, w, gO))
		case "Sprites":
			var sp models.Sprite = models.LoadSprites(p.Sprites)[edStatus.selected.id]
			edStatus.editorPanel.Add(buildGameSpriteEditor(p, w, sp))
		case "Rooms":
			var roo models.Room = models.LoadRooms(p.Rooms)[edStatus.selected.id]
			edStatus.editorPanel.Add(buildRoomEditor(p, w, roo))
		}
		edStatus.editorPanel.Refresh()
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
				sp.Directory = filepath.Join(p.Directory, "Sprites", sp.Name)
				sp.Columns = 1
				utils.EnsureDir(sp.Directory)
				p.Sprites = append(p.Sprites, filepath.Join(sp.Directory, sp.Name+".lsp"))

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
