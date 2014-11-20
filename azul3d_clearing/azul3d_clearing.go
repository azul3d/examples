// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates clearing the window.
package main

import (
	"image"

	"azul3d.org/gfx.v2-dev"
	"azul3d.org/gfx.v2-dev/window"
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, r gfx.Renderer) {
	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})

		// Clear a few rectangles on the window using different background
		// colors.
		r.Clear(image.Rect(0, 50, 800, 400), gfx.Color{0, 1, 0, 1})
		r.Clear(image.Rect(50, 50, 750, 400), gfx.Color{1, 0, 0, 1})
		r.Clear(image.Rect(50, 100, 750, 350), gfx.Color{0, 0.5, 0.5, 1})
		r.Clear(image.Rect(100, 150, 700, 300), gfx.Color{1, 1, 0, 1})

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
