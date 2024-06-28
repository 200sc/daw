package oak

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"time"
)

// ScreenShot takes a snap shot of the window's image content.
// ScreenShot is not safe to call while an existing ScreenShot call has
// yet to finish executing. This could change in the future.
func (w *Window) ScreenShot() *image.RGBA {
	shotCh := make(chan *image.RGBA)
	// We need to take the shot when the screen is not being redrawn
	// We know the screen has everything drawn on it when it is published
	w.prePublish = func(rgba *image.RGBA) {
		// Copy the buffer
		bds := rgba.Bounds()
		copy := image.NewRGBA(bds)
		for x := bds.Min.X; x < bds.Max.X; x++ {
			for y := bds.Min.Y; y < bds.Max.Y; y++ {
				copy.Set(x, y, rgba.RGBAAt(x, y))
			}
		}
		shotCh <- copy
	}
	out := <-shotCh
	w.ClearScreenFilter()
	return out
}

// gifShot is internally used by RecordGIF
func (w *Window) gifShot() *image.Paletted {
	shotCh := make(chan *image.Paletted)
	// We need to take the shot when the screen is not being redrawn
	// We know the screen has everything drawn on it when it is published
	w.prePublish = func(rgba *image.RGBA) {
		// Copy the buffer
		bds := rgba.Bounds()
		copy := image.NewPaletted(bds, palette.Plan9)
		draw.Draw(copy, bds, rgba, zeroPoint, draw.Src)
		shotCh <- copy
	}
	out := <-shotCh
	w.ClearScreenFilter()
	return out
}

// RecordGIF will start recording frames via screen shots with the given
// time delay (in 1/100ths of a second) between frames. When the returned
// stop function is called, the frames will be compiled into a gif.
func (w *Window) RecordGIF(hundredths int) (stop func() *gif.GIF) {
	cancel := make(chan struct{})
	out := make(chan *gif.GIF)
	delay := time.Duration(hundredths) * time.Millisecond * 10
	go func() {
		g := &gif.GIF{}
		for {
			select {
			case <-time.After(delay):
			case <-cancel:
				out <- g
				return
			}
			shot := w.gifShot()
			g.Image = append(g.Image, shot)
			g.Delay = append(g.Delay, hundredths)
		}
	}()
	return func() *gif.GIF {
		cancel <- struct{}{}
		return <-out
	}
}
