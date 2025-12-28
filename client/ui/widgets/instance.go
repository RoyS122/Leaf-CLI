package widgets

import (
	"fmt"
	"image"
	"leafcli/models"
	"leafcli/utils"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type InstanceWidget struct {
	widget.BaseWidget

	Instance *models.Instance
	Image    *canvas.Image

	position fyne.Position
	Selected bool

	OnSelected           func(*models.Instance)
	OnRoomUpdateinstance func()
}

func NewInstanceWidget(instance *models.Instance) *InstanceWidget {
	fmt.Println("Creating InstanceWidget for instance at", instance.X, instance.Y)
	w := &InstanceWidget{
		Instance: instance,
		position: fyne.NewPos(float32(instance.X), float32(instance.Y)),
	}

	// Charger le GameObject
	obj := instance.GetGameObject()
	if obj.Name == "" {
		fmt.Println("Failed to load GameObject for instance")
		return nil
	}

	sprite, err := models.LoadSpriteFromFile(obj.Sprite)
	if err != nil {
		fmt.Println("Failed to load Sprite for GameObject:", obj.Name)
		return nil
	}

	// Charger l’image
	file, err := os.Open(sprite.ImagePath)
	if err != nil {
		fmt.Println("Error opening sprite image:", sprite.ImagePath)
		return nil
	}
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("fail at decoding image:", err)
		return nil
	}

	// Image canvas
	img := canvas.NewImageFromImage(srcImg)
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels
	img.Resize(fyne.NewSize(64, 64))

	// Sprite sheet
	if sprite.Columns > 0 {
		frames := utils.SplitImage(srcImg, 1, int(sprite.Columns))
		utils.ChangeImage(&img, frames, 0)
	}

	img.Move(fyne.NewPos(0, 0))
	img.Refresh()

	w.Image = img
	w.ExtendBaseWidget(w)
	w.Resize(img.Size()) // ❗ OBLIGATOIRE
	w.Move(w.position)
	fmt.Println("Created InstanceWidget for GameObject:", obj.Name)
	return w
}
func (w *InstanceWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.Image)
}

func (w *InstanceWidget) Tapped(*fyne.PointEvent) {
	w.Selected = !w.Selected

	if w.OnSelected != nil {
		w.OnSelected(w.Instance)
	}

	// Feedback visuel simple (optionnel)
	if w.Selected {
		w.Image.Translucency = 0.8

	} else {
		w.Image.Translucency = 0

	}
	canvas.Refresh(w.Image)
}

func (w *InstanceWidget) Dragged(ev *fyne.DragEvent) {
	w.position = w.position.Add(ev.Dragged)

	w.Move(w.position)
	w.Image.Translucency = 0.8
	// Synchronise le modèle
	w.Instance.X = int(w.position.X)
	w.Instance.Y = int(w.position.Y)
	// w.OnRoomUpdateinstance()
	canvas.Refresh(w.Image)
}

func (w *InstanceWidget) DragEnd() {
	w.Image.Translucency = 0
	w.OnRoomUpdateinstance()
	canvas.Refresh(w.Image)
}
