// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates clearing the window.
package main

import (
	"azul3d.org/chippy.v1"
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v1"
	"image"
)

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})

		// Clear a few rectangles on the window using different background
		// colors.
		r.Clear(image.Rect(0, 100, 720, 380), gfx.Color{0, 1, 0, 1})
		r.Clear(image.Rect(100, 100, 620, 380), gfx.Color{1, 0, 0, 1})
		r.Clear(image.Rect(100, 200, 620, 280), gfx.Color{0, 0.5, 0.5, 1})
		r.Clear(image.Rect(200, 200, 520, 280), gfx.Color{1, 1, 0, 1})

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop)
}
