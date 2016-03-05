// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates texture coordinates.
package main

import (
	_ "image/png"
	"log"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/gfxutil"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/keyboard"
	"azul3d.org/engine/lmath"

	"azul3d.org/examples/abs"
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, d gfx.Device) {
	// Create a new orthographic (2D) camera.
	cam := camera.NewOrtho(d.Bounds())

	// Move the camera back two units away from the card.
	cam.SetPos(lmath.Vec3{0, -2, 0})

	// Read the GLSL shaders from disk.
	shader, err := gfxutil.OpenShader(abs.Path("azul3d_texcoords/texcoords"))
	if err != nil {
		log.Fatal(err)
	}

	// Open the texture.
	tex, err := gfxutil.OpenTexture(abs.Path("azul3d_texcoords/texture_coords_1024x1024.png"))
	if err != nil {
		log.Fatal(err)
	}
	tex.Format = gfx.DXT1RGBA

	// Create a card mesh.
	cardMesh := gfx.NewMesh()
	cardMesh.Vertices = []gfx.Vec3{
		// Bottom-left triangle.
		{-1, 0, -1},
		{1, 0, -1},
		{-1, 0, 1},

		// Top-right triangle.
		{-1, 0, 1},
		{1, 0, -1},
		{1, 0, 1},
	}
	cardMesh.TexCoords = []gfx.TexCoordSet{
		{
			Slice: []gfx.TexCoord{
				{0, 1},
				{1, 1},
				{0, 0},

				{0, 0},
				{1, 1},
				{1, 0},
			},
		},
	}

	// Create a card object.
	card := gfx.NewObject()
	card.State = gfx.NewState()
	card.AlphaMode = gfx.AlphaToCoverage
	card.Shader = shader
	card.Textures = []*gfx.Texture{tex}
	card.Meshes = []*gfx.Mesh{cardMesh}

	// Create a channel of events.
	events := make(chan window.Event, 256)

	// Have the window notify our channel whenever events occur.
	w.Notify(events, window.FramebufferResizedEvents|window.KeyboardTypedEvents)

	for {
		// Handle each pending event.
		window.Poll(events, func(e window.Event) {
			switch ev := e.(type) {
			case window.FramebufferResized:
				// Update the camera's projection matrix for the new width and
				// height.
				cam.Update(d.Bounds())

			case keyboard.Typed:
				if ev.S == "m" || ev.S == "M" {
					// Toggle mipmapping on the texture.
					if tex.MinFilter == gfx.LinearMipmapLinear {
						tex.MinFilter = gfx.Linear
					} else {
						tex.MinFilter = gfx.LinearMipmapLinear
					}
				}
			}
		})

		// Center the card in the window.
		b := d.Bounds()
		card.SetPos(lmath.Vec3{float64(b.Dx()) / 2.0, 0, float64(b.Dy()) / 2.0})

		// Scale the card to fit the window.
		s := float64(b.Dy()) / 2.0 // Card is two units wide, so divide by two.
		card.SetScale(lmath.Vec3{s, s, s})

		// Clear color and depth buffers.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
		d.ClearDepth(d.Bounds(), 1.0)

		// Draw the textured card.
		d.Draw(d.Bounds(), card, cam)

		// Render the whole frame.
		d.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
