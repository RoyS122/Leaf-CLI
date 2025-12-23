package widgets

import (
	"image/color"
	"leafcli/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildRoomViewport(r models.Room) fyne.CanvasObject {
	const (
		roomWidth  = 500
		roomHeight = 300
		gridSize   = 32
	)

	// Fond de la room
	background := canvas.NewRectangle(color.NRGBA{R: 40, G: 40, B: 40, A: 255})
	background.Resize(fyne.NewSize(roomWidth, roomHeight))

	objects := []fyne.CanvasObject{
		background,
	}

	// Couleur de la grille (légèrement transparente)
	gridColor := color.NRGBA{R: 255, G: 255, B: 255, A: 40}

	// Lignes verticales
	for x := gridSize; x < roomWidth; x += gridSize {
		line := canvas.NewLine(gridColor)
		line.Position1 = fyne.NewPos(float32(x), 0)
		line.Position2 = fyne.NewPos(float32(x), roomHeight)
		objects = append(objects, line)
	}

	// Lignes horizontales
	for y := gridSize; y < roomHeight; y += gridSize {
		line := canvas.NewLine(gridColor)
		line.Position1 = fyne.NewPos(0, float32(y))
		line.Position2 = fyne.NewPos(roomWidth, float32(y))
		objects = append(objects, line)
	}

	// Conteneur libre
	roomContent := container.NewWithoutLayout(objects...)
	roomContent.Resize(fyne.NewSize(roomWidth, roomHeight))

	// Padding extérieur (optionnel)
	return container.NewPadded(roomContent)
}

type RoomViewport struct {
	widget.BaseWidget
	content fyne.CanvasObject
}

func NewRoomViewport(room models.Room) *RoomViewport {
	rv := &RoomViewport{
		content: BuildRoomViewport(room),
	}
	rv.ExtendBaseWidget(rv)
	return rv
}

func (r *RoomViewport) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.content)
}

func (r *RoomViewport) MinSize() fyne.Size {
	return fyne.NewSize(500, 300)
}
