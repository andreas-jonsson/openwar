/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

import (
	"image"
	"image/color"
	"image/draw"
	"unsafe"

	"github.com/andreas-jonsson/go-sdl2/sdl"
)

type sdlRenderer struct {
	window           *sdl.Window
	backBuffer       *image.RGBA
	internalHWBuffer *sdl.Texture
	internalRenderer *sdl.Renderer
}

func NewRenderer(w, h int, data ...interface{}) (Renderer, error) {
	var (
		r   sdlRenderer
		err error
	)

	r.window, err = sdl.CreateWindow(data[0].(string), sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	r.backBuffer = image.NewRGBA(image.Rect(0, 0, 320, 200))

	renderer, err := sdl.CreateRenderer(r.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	r.internalRenderer = renderer

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")
	renderer.SetLogicalSize(320, 240)
	renderer.SetDrawColor(0, 0, 0, 255)

	r.internalHWBuffer, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, 320, 200)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (p *sdlRenderer) ToggleFullscreen() {
	isFullscreen := (p.window.GetFlags() & sdl.WINDOW_FULLSCREEN) != 0
	if isFullscreen {
		p.window.SetFullscreen(0)
		sdl.ShowCursor(1)
	} else {
		p.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
		sdl.ShowCursor(0)
	}
}

func (p *sdlRenderer) Clear() {
	p.internalRenderer.Clear()
	draw.Draw(p.backBuffer, p.backBuffer.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
}

func (p *sdlRenderer) Present() {
	p.internalHWBuffer.Update(nil, unsafe.Pointer(&p.backBuffer.Pix[0]), p.backBuffer.Stride)
	p.internalRenderer.Copy(p.internalHWBuffer, nil, nil)
	p.internalRenderer.Present()
}

func (p *sdlRenderer) Shutdown() {
	p.window.Destroy()
	p.internalHWBuffer.Destroy()
	p.internalRenderer.Destroy()
}

func (p *sdlRenderer) BackBuffer() draw.Image {
	return p.backBuffer
}

func (p *sdlRenderer) Blit(src image.Image, sp image.Point) {
	draw.Draw(p.backBuffer, p.backBuffer.Bounds(), src, sp, draw.Src)
}

func (p *sdlRenderer) BlitPal(src *image.Paletted, pal color.Palette, sp image.Point) {
	orgPal := src.Palette
	src.Palette = pal

	draw.Draw(p.backBuffer, p.backBuffer.Bounds(), src, sp, draw.Src)

	src.Palette = orgPal
}
