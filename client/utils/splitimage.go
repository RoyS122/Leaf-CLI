package utils

import (
	"image"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func SplitImage(srcIMG image.Image, rows, cols int) (image_splited []*image.RGBA) {
	if rows <= 0 || cols <= 0 {
		return
	}
	bounds := srcIMG.Bounds()
	frameWidth := bounds.Dx() / cols
	frameHeight := bounds.Dy() / rows
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			rect := image.Rect(0, 0, frameWidth, frameHeight)
			croppedImg := image.NewRGBA(rect)
			// si srcIMG.Bounds().Min != (0,0), copier depuis la bonne origine
			srcPt := image.Point{X: srcIMG.Bounds().Min.X + x*frameWidth, Y: srcIMG.Bounds().Min.Y + y*frameHeight}
			draw.Draw(croppedImg, rect, srcIMG, srcPt, draw.Src)
			image_splited = append(image_splited, croppedImg)
		}
	}
	return image_splited
}

func ChangeImage(imgcanvas **canvas.Image, frames []*image.RGBA, index int) {
	*imgcanvas = canvas.NewImageFromImage(frames[index])
	(*imgcanvas).SetMinSize(fyne.NewSize(float32(frames[index].Bounds().Dx()), 64))
	(*imgcanvas).FillMode = canvas.ImageFillContain
	(*imgcanvas).ScaleMode = canvas.ImageScalePixels
	(*imgcanvas).Resize(fyne.NewSize(float32(frames[index].Bounds().Dx()), 64))
}
