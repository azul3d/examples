// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays a TMX tiled map.
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/keyboard"
	"azul3d.org/engine/lmath"
	"azul3d.org/engine/mouse"
	"azul3d.org/engine/tmx"

	"azul3d.org/examples/abs"
)

// setOrthoScale sets the camera's projection matrix to an orthographic one
// using the given viewing rectangle. It performs scaling with the viewing
// rectangle.
func setOrthoScale(c *camera.Camera, view image.Rectangle, scale float64) {
	w := float64(c.View.Dx())
	w *= scale
	w = float64(int((w / 2.0)))

	h := float64(c.View.Dy())
	h *= scale
	h = float64(int((h / 2.0)))

	m := lmath.Mat4Ortho(-w, w, -h, h, c.Near, c.Far)
	c.P = gfx.ConvertMat4(m)
}

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, d gfx.Device) {
	// Create a new orthographic (2D) camera.
	cam := camera.NewOrtho(d.Bounds())
	camZoom := 1.0       // 1x zoom
	camZoomSpeed := 0.01 // 0.01x zoom for each scroll wheel click.
	camMinZoom := 0.1

	// updateCamera simply calls setOrthoScale with the values above.
	updateCamera := func() {
		if camZoom < camMinZoom {
			camZoom = camMinZoom
		}
		setOrthoScale(cam, d.Bounds(), camZoom)
	}

	// Update the camera now.
	updateCamera()

	// Move the camera back two units away from the card.
	cam.SetPos(lmath.Vec3{0, -2, 0})

	// Load TMX map file.
	tmxMap, layers, err := tmx.LoadFile(*mapFile, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create an event mask for the events we are interested in.
	evMask := window.FramebufferResizedEvents
	evMask |= window.CursorMovedEvents
	evMask |= window.MouseEvents
	evMask |= window.MouseScrolledEvents
	evMask |= window.KeyboardTypedEvents

	// Create a channel of events.
	events := make(chan window.Event, 256)

	// Have the window notify our channel whenever events occur.
	w.Notify(events, evMask)

	handleEvent := func(e window.Event) {
		switch ev := e.(type) {
		case window.FramebufferResized:
			// Update the camera's to account for the new width and height.
			updateCamera()

		case mouse.ButtonEvent:
			if ev.Button == mouse.Left && ev.State == mouse.Up {
				// Toggle mouse grab.
				props := w.Props()
				props.SetCursorGrabbed(!props.CursorGrabbed())
				w.Request(props)
			}

		case mouse.Scrolled:
			// Zoom and update the camera.
			camZoom -= ev.Y * camZoomSpeed
			updateCamera()

		case window.CursorMoved:
			if ev.Delta {
				p := lmath.Vec3{ev.X, 0, -ev.Y}
				p = p.MulScalar(camZoom)
				cam.SetPos(cam.Pos().Add(p))
			}

		case keyboard.Typed:
			switch ev.S {
			case "m":
				// Toggle MSAA now.
				msaa := !d.MSAA()
				d.SetMSAA(msaa)
				fmt.Println("MSAA Enabled?", msaa)
			case "r":
				cam.SetPos(lmath.Vec3{0, -2, 0})
			}
		}
	}

	for {
		// Handle events.
		window.Poll(events, handleEvent)

		// Clear color and depth buffers.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
		d.ClearDepth(d.Bounds(), 1.0)

		// Draw the TMX map to the screen.
		for _, layer := range tmxMap.Layers {
			objects, ok := layers[layer.Name]
			if ok {
				for _, obj := range objects {
					d.Draw(d.Bounds(), obj, cam)
				}
			}
		}

		// Render the whole frame.
		d.Render()
	}
}

var (
	defaultMapFile = abs.Path("azul3d_tmx/data/test_base64.tmx")
	mapFile        = flag.String("file", defaultMapFile, "tmx map file to load")
)

func init() {
	flag.Parse()
}

func main() {
	window.Run(gfxLoop, nil)
}
