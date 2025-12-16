package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Forcer le backend OpenGL au lieu de GTK
	os.Setenv("FYNE_FONT", "/mingw64/share/fonts/TTF/DejaVuSans.ttf")
	os.Setenv("FYNE_THEME", "light")
	
	a := app.NewWithID("com.example.hello")
	w := a.NewWindow("Hello")
	
	hello := widget.NewLabel("Hello Fyne!")
	button := widget.NewButton("Click me!", func() {
		hello.SetText("Welcome :)")
		fmt.Println("Button clicked!")
	})
	
	w.SetContent(container.NewVBox(
		hello,
		button,
	))
	
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}