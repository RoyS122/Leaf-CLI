package widgets

import (
	"image/color"
	"leafcli/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// RoomViewport gère l'affichage d'une Room et de ses instances
type RoomViewport struct {
	widget.BaseWidget
	room *models.Room
	root *fyne.Container // racine stable
}

// NewRoomViewport crée un RoomViewport pour la room donnée
func NewRoomViewport(room *models.Room) *RoomViewport {
	root := container.NewWithoutLayout()
	rv := &RoomViewport{
		room: room,
		root: root,
	}
	rv.ExtendBaseWidget(rv)
	rv.Reload() // remplir le root
	return rv
}

// Reload reconstruit le contenu du viewport
func (r *RoomViewport) Reload() {
	const (
		roomWidth  = 512
		roomHeight = 288
		gridSize   = 32
	)

	objects := []fyne.CanvasObject{}

	// Fond
	bg := canvas.NewRectangle(color.NRGBA{R: 40, G: 40, B: 40, A: 255})
	bg.Resize(fyne.NewSize(roomWidth, roomHeight))
	objects = append(objects, bg)

	// Grille
	gridColor := color.NRGBA{R: 255, G: 255, B: 255, A: 40}
	for x := gridSize; x < roomWidth; x += gridSize {
		line := canvas.NewLine(gridColor)
		line.Position1 = fyne.NewPos(float32(x), 0)
		line.Position2 = fyne.NewPos(float32(x), roomHeight)
		objects = append(objects, line)
	}
	for y := gridSize; y < roomHeight; y += gridSize {
		line := canvas.NewLine(gridColor)
		line.Position1 = fyne.NewPos(0, float32(y))
		line.Position2 = fyne.NewPos(roomWidth, float32(y))
		objects = append(objects, line)
	}

	// Instances
	for i := range r.room.GOInstances {
		inst := &r.room.GOInstances[i]
		w := NewInstanceWidget(inst)
		if w != nil {
			w.OnSelected = func(i *models.Instance) {}
			w.OnRoomUpdateinstance = func() { r.room.Save() }
			objects = append(objects, w)
		}
	}

	// Mettre à jour le contenu du root existant
	r.root.Objects = objects
	r.root.Refresh()
	r.BaseWidget.Refresh()
}

// Implémentation du renderer
func (r *RoomViewport) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.root)
}

func (r *RoomViewport) MinSize() fyne.Size {
	return fyne.NewSize(500, 300)
}
