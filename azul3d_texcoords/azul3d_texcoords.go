// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates texture coordinates.
package main

import (
	"go/build"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"azul3d.org/gfx.v2-dev"
	"azul3d.org/gfx.v2-dev/window"
	"azul3d.org/keyboard.v1"
	"azul3d.org/lmath.v1"
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
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camNear := 0.01
	camFar := 1000.0
	camera.SetOrtho(d.Bounds(), camNear, camFar)

	// Move the camera back two units away from the card.
	camera.SetPos(lmath.Vec3{0, -2, 0})

	// Loading shader files
	glslVert, err := ioutil.ReadFile(absPath("azul3d_texcoords/texcoords.vert"))
	if err != nil {
		log.Fatal(err)
	}
	glslFrag, err := ioutil.ReadFile(absPath("azul3d_texcoords/texcoords.frag"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSL = &gfx.GLSLSources{
		Vertex:   glslVert,
		Fragment: glslFrag,
	}

	// Load the picture.
	f, err := os.Open(absPath("azul3d_texcoords/texture_coords_1024x1024.png"))
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Create new texture.
	tex := gfx.NewTexture()
	tex.Source = img
	tex.MinFilter = gfx.LinearMipmapLinear
	tex.MagFilter = gfx.Linear
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
	card.AlphaMode = gfx.AlphaToCoverage
	card.Shader = shader
	card.Textures = []*gfx.Texture{tex}
	card.Meshes = []*gfx.Mesh{cardMesh}

	go func() {
		// Create a channel of events.
		events := make(chan window.Event, 256)

		// Have the window notify our channel whenever events occur.
		w.Notify(events, window.FramebufferResizedEvents|window.KeyboardTypedEvents)

		for e := range events {
			switch ev := e.(type) {
			case window.FramebufferResized:
				// Update the camera's projection matrix for the new width and
				// height.
				camera.Lock()
				camera.SetOrtho(d.Bounds(), camNear, camFar)
				camera.Unlock()

			case keyboard.TypedEvent:
				if ev.Rune == 'm' || ev.Rune == 'M' {
					// Toggle mipmapping on the texture.
					tex.Lock()
					if tex.MinFilter == gfx.LinearMipmapLinear {
						tex.MinFilter = gfx.Linear
					} else {
						tex.MinFilter = gfx.LinearMipmapLinear
					}
					tex.Unlock()
				}
			}
		}
	}()

	for {
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
		d.Draw(d.Bounds(), card, camera)

		// Render the whole frame.
		d.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
