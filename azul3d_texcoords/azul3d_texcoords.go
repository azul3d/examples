// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates texture coordinates.
package main

import (
	"azul3d.org/chippy.v1"
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v1"
	math "azul3d.org/lmath.v1"
	"image"
	"log"
	"os"
)

var glslVert = []byte(`
#version 120

attribute vec3 Vertex;
attribute vec2 TexCoord0;

uniform mat4 MVP;

varying vec2 tc0;

void main()
{
	tc0 = TexCoord0;
	gl_Position = MVP * vec4(Vertex, 1.0);
}
`)

var glslFrag = []byte(`
#version 120

varying vec2 tc0;

uniform sampler2D Texture0;
uniform bool BinaryAlpha;

void main()
{
	gl_FragColor = texture2D(Texture0, tc0);
	if(BinaryAlpha && gl_FragColor.a < 0.5) {
		discard;
	}
}
`)

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camNear := 0.01
	camFar := 1000.0
	camera.SetOrtho(r.Bounds(), camNear, camFar)

	// Move the camera back two units away from the card.
	camera.SetPos(math.Vec3{0, -2, 0})

	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSLVert = glslVert
	shader.GLSLFrag = glslFrag

	// Load the picture.
	f, err := os.Open("src/azul3d.org/examples.v1/assets/textures/texture_coords_1024x1024.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Create new texture.
	tex := &gfx.Texture{
		Source:    img,
		MinFilter: gfx.LinearMipmapLinear,
		MagFilter: gfx.Linear,
		Format:    gfx.DXT1RGBA,
	}

	// Create a card mesh.
	cardMesh := &gfx.Mesh{
		Vertices: []gfx.Vec3{
			// Bottom-left triangle.
			{-1, 0, -1},
			{1, 0, -1},
			{-1, 0, 1},

			// Top-right triangle.
			{-1, 0, 1},
			{1, 0, -1},
			{1, 0, 1},
		},
		TexCoords: []gfx.TexCoordSet{
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
		},
	}

	// Create a card object.
	card := gfx.NewObject()
	card.AlphaMode = gfx.AlphaToCoverage
	card.Shader = shader
	card.Textures = []*gfx.Texture{tex}
	card.Meshes = []*gfx.Mesh{cardMesh}

	go func() {
		for e := range w.Events() {
			switch e.(type) {
			case chippy.ResizedEvent:
				// Update the camera's projection matrix for the new width and
				// height.
				camera.Lock()
				camera.SetOrtho(r.Bounds(), camNear, camFar)
				camera.Unlock()
			}
		}
	}()

	for {
		// Center the card in the window.
		b := r.Bounds()
		card.SetPos(math.Vec3{float64(b.Dx()) / 2.0, 0, float64(b.Dy()) / 2.0})

		// Scale the card to fit the window.
		s := float64(b.Dy()) / 2.0 // Card is two units wide, so divide by two.
		card.SetScale(math.Vec3{s, s, s})

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Draw the textured card.
		r.Draw(image.Rect(0, 0, 0, 0), card, camera)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop)
}
