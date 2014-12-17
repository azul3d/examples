// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Generates a mandelbrot on the CPU and displays it with the GPU.
package main

import (
	"fmt"
	"go/build"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	gmath "math"
	"os"
	"path/filepath"

	"azul3d.org/gfx.v2-dev"
	"azul3d.org/gfx.v2-dev/window"
	"azul3d.org/keyboard.v1"
	"azul3d.org/mouse.v1"
)

// This helper function is not an important example concept, please ignore it.
//
// absPath the absolute path to an file given one relative to the examples
// directory:
//  $GOPATH/src/azul3d.org/examples.v1
var examplesDir string

func absPath(relPath string) string {
	if len(examplesDir) == 0 {
		// Find assets directory.
		for _, path := range filepath.SplitList(build.Default.GOPATH) {
			path = filepath.Join(path, "src/azul3d.org/examples.v1")
			if _, err := os.Stat(path); err == nil {
				examplesDir = path
				break
			}
		}
	}
	return filepath.Join(examplesDir, relPath)
}

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, d gfx.Device) {
	// Loading shader files
	glslVert, err := ioutil.ReadFile(absPath("azul3d_mandel/mandel.vert"))
	if err != nil {
		log.Fatal(err)
	}
	glslFrag, err := ioutil.ReadFile(absPath("azul3d_mandel/mandel.frag"))
	if err != nil {
		log.Fatal(err)
	}
	
	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSL = &gfx.GLSLSources{
		Vertex:   glslVert,
		Fragment: glslFrag,
	}

	// Create a card mesh.
	cardMesh := gfx.NewMesh()
	cardMesh.Vertices = []gfx.Vec3{
		// Left triangle.
		{-1, 1, 0},  // Left-Top
		{-1, -1, 0}, // Left-Bottom
		{1, -1, 0},  // Right-Bottom

		// Right triangle.
		{-1, 1, 0}, // Left-Top
		{1, -1, 0}, // Right-Bottom
		{1, 1, 0},  // Right-Top
	}
	cardMesh.TexCoords = []gfx.TexCoordSet{
		{
			Slice: []gfx.TexCoord{
				// Left triangle.
				{0, 0},
				{0, 1},
				{1, 1},

				// Right triangle.
				{0, 0},
				{1, 1},
				{1, 0},
			},
		},
	}

	// Create a card object.
	card := gfx.NewObject()
	card.Shader = shader
	card.Textures = []*gfx.Texture{nil}
	card.Meshes = []*gfx.Mesh{cardMesh}

	// Create a texture.
	zoom := 1.0
	x := -0.5
	y := 0.0
	res := 8
	maxIter := 1000
	updateTex := func() {
		width, height := d.Bounds().Dx(), d.Bounds().Dy()
		mbrot := Mandelbrot(width/res, height/res, maxIter, zoom, x, y)

		// Insert a small red square in the top-left of the image for ensuring
		// proper orientation exists in textures (this is just for testing).
		for x := 0; x < 20; x++ {
			for y := 0; y < 20; y++ {
				mbrot.Set(x, y, color.RGBA{255, 0, 0, 255})
			}
		}

		// Create new texture and ask the renderer to load it. We don't use DXT
		// compression because those textures cannot be downloaded.
		tex := gfx.NewTexture()
		tex.Source = mbrot
		tex.MinFilter = gfx.Nearest
		tex.MagFilter = gfx.Nearest

		onLoad := make(chan *gfx.Texture, 1)
		d.LoadTexture(tex, onLoad)
		<-onLoad

		// Swap the texture with the old one on the card.
		card.Lock()
		card.Textures[0] = tex
		card.Unlock()
	}
	updateTex()

	go func() {
		handle := func(e window.Event) (needUpdate bool) {
			switch ev := e.(type) {
			case mouse.Event:
				if ev.Button == mouse.Left && ev.State == mouse.Down {
					props := w.Props()
					props.SetCursorGrabbed(!props.CursorGrabbed())
					w.Request(props)
				}

				if ev.Button == mouse.Right && ev.State == mouse.Down {
					res += 2
					if res > 8 {
						res = 4
					}
					return true
				}

			case mouse.Scrolled:
				zoom += ev.Y * 0.06 * gmath.Abs(zoom)

			case keyboard.TypedEvent:
				if ev.Rune == 's' || ev.Rune == 'S' {
					fmt.Println("Writing texture to file...")
					// Download the image from the graphics hardware and save
					// it to disk.
					complete := make(chan image.Image, 1)

					// Lock the card/texture.
					card.RLock()
					card.Textures[0].Lock()

					// Begin downloading it's texture.
					card.Textures[0].Download(image.Rect(0, 0, 640, 480), complete)

					// Unlock the texture/card.
					card.Textures[0].Unlock()
					card.RUnlock()

					img := <-complete // Wait for download to complete.
					if img == nil {
						fmt.Println("Failed to download texture.")
					} else {
						// Save to png.
						f, err := os.Create("mandel.png")
						if err != nil {
							log.Fatal(err)
						}
						err = png.Encode(f, img)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("Wrote texture to mandel.png")
					}
				}

			case window.CursorMoved:
				if ev.Delta {
					x += (ev.X / 900.0) / gmath.Abs(zoom)
					y += (ev.Y / 900.0) / gmath.Abs(zoom)
					return true
				}
			}
			return false
		}

		// Create an event mask for the events we are interested in.
		evMask := window.MouseEvents
		evMask |= window.MouseScrolledEvents
		evMask |= window.KeyboardTypedEvents
		evMask |= window.CursorMovedEvents

		// Create a channel of events.
		events := make(chan window.Event, 256)

		// Have the window notify our channel whenever events occur.
		w.Notify(events, evMask)

		// Wait for events, we process them in large batches because updateTex
		// calculate a mandelbrot on the CPU and it's very slow.
		for {
			e := <-events
			needUpdate := handle(e)
			l := len(events)
			for i := 0; i < l; i++ {
				if handle(<-events) {
					needUpdate = true
				}
			}
			if needUpdate {
				// Generate new mandel texture.
				updateTex()
			}
		}
	}()

	for {
		// Clear color and depth buffers.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
		d.ClearDepth(d.Bounds(), 1.0)

		// Draw the card to the screen.
		d.Draw(d.Bounds(), card, nil)

		// Render the whole frame.
		d.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
