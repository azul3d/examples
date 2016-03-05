// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates clearing multiple windows.
package main

import (
	"image"
	"log"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/window"
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, d gfx.Device) {
	for {
		// Clear the entire area.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})

		// Clear a few rectangles on the window using different background
		// colors.
		d.Clear(image.Rect(0, 50, 800, 400), gfx.Color{0, 1, 0, 1})
		d.Clear(image.Rect(50, 50, 750, 400), gfx.Color{1, 0, 0, 1})
		d.Clear(image.Rect(50, 100, 750, 350), gfx.Color{0, 0.5, 0.5, 1})
		d.Clear(image.Rect(100, 150, 700, 300), gfx.Color{1, 1, 0, 1})

		// Render the whole frame.
		d.Render()
	}
}

func main() {
	// Create our windows in a seperate goroutine. A seperate goroutine is
	// needed because New cannot complete unless the MainLoop is running.
	go func() {
		// Create our windows.
		for i := 0; i < 2; i++ {
			props := window.NewProps()
			props.SetTitle("Window {FPS}")
			w, r, err := window.New(props)
			if err != nil {
				log.Fatal(err)
			}

			// Start a graphics loop for the window.
			go gfxLoop(w, r)
		}
	}()

	// Enter the main loop.
	window.MainLoop()
}
