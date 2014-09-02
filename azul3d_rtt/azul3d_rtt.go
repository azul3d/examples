// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates texture coordinates.
package main

import (
	"image"
	"log"

	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/keyboard.v1"
	"azul3d.org/lmath.v1"
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

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, r gfx.Renderer) {
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camFOV := 75.0
	camNear := 0.0001
	camFar := 1000.0
	camera.SetPersp(r.Bounds(), camFOV, camNear, camFar)

	// Move the camera back two units away from the card.
	camera.SetPos(lmath.Vec3{0, -2, 0})

	// Create a texture to hold the color data of our render-to-texture.
	rtColor := gfx.NewTexture()
	rtColor.MinFilter = gfx.LinearMipmapLinear
	rtColor.MagFilter = gfx.Linear

	// Choose a render to texture format.
	cfg := r.GPUInfo().RTTFormats.ChooseConfig(gfx.Precision{
		// We want 24/bpp RGB color buffer.
		RedBits: 8, GreenBits: 8, BlueBits: 8,

		// We could also request a depth or stencil buffer here, by simply
		// using the lines:
		// DepthBits: 24,
		// StencilBits: 24,
	}, true)

	// Print the configuration we chose.
	log.Printf("RTT ColorFormat=%v, DepthFormat=%v, StencilFormat=%v\n", cfg.ColorFormat, cfg.DepthFormat, cfg.StencilFormat)

	// Color buffer will go into our rtColor texture.
	cfg.Color = rtColor

	// We will render to a 512x512 area.
	cfg.Bounds = image.Rect(0, 0, 512, 512)

	// Create our render-to-texture canvas.
	rtCanvas := r.RenderToTexture(cfg)
	if rtCanvas == nil {
		// Important! Check if the canvas is nil. If it is their graphics
		// hardware doesn't support render to texture. Sorry!
		log.Fatal("Graphics hardware does not support render to texture.")
	}

	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSLVert = glslVert
	shader.GLSLFrag = glslFrag

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
	card.FaceCulling = gfx.NoFaceCulling
	card.AlphaMode = gfx.AlphaToCoverage
	card.Shader = shader
	card.Textures = []*gfx.Texture{rtColor}
	card.Meshes = []*gfx.Mesh{cardMesh}

	// Spawn a goroutine to handle events.
	go func() {
		// Create an event mask for the events we are interested in.
		evMask := window.FramebufferResizedEvents
		evMask |= window.KeyboardTypedEvents

		// Create a channel of events.
		events := make(chan window.Event, 256)

		// Have the window notify our channel whenever events occur.
		w.Notify(events, evMask)

		for e := range events {
			switch ev := e.(type) {
			case window.FramebufferResized:
				// Update the camera's projection matrix for the new width and
				// height.
				camera.Lock()
				camera.SetPersp(r.Bounds(), camFOV, camNear, camFar)
				camera.Unlock()

			case keyboard.TypedEvent:
				if ev.Rune == 'm' || ev.Rune == 'M' {
					// Toggle mipmapping.
					rtColor.Lock()
					if rtColor.MinFilter == gfx.LinearMipmapLinear {
						rtColor.MinFilter = gfx.Linear
					} else {
						rtColor.MinFilter = gfx.LinearMipmapLinear
					}
					rtColor.Unlock()
				}
			}
		}
	}()

	// Draw some colored stripes onto the render to texture canvas. The result
	// is stored in the rtColor texture, and we can then display it on a card
	// below without even rendering the stripes every frame.
	stripeColor1 := gfx.Color{1, 0, 0, 1} // red
	stripeColor2 := gfx.Color{0, 1, 0, 1} // green
	stripeWidth := 12                     // pixels
	flipColor := false
	b := rtCanvas.Bounds()
	for i := 0; (i * stripeWidth) < b.Dx(); i++ {
		flipColor = !flipColor
		x := i * stripeWidth
		dst := image.Rect(x, b.Min.Y, x+stripeWidth, b.Max.Y)
		if flipColor {
			rtCanvas.Clear(dst, stripeColor1)
		} else {
			rtCanvas.Clear(dst, stripeColor2)
		}
	}

	// Render the rtCanvas to the rtColor texture.
	rtCanvas.Render()

	for {
		// Rotate the card on the Z axis 15 degrees/sec.
		rot := card.Rot()
		card.SetRot(lmath.Vec3{
			X: rot.X,
			Y: rot.Y,
			Z: rot.Z + (15 * r.Clock().Dt()),
		})

		// Clear color and depth buffers.
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Draw the card.
		r.Draw(image.Rect(0, 0, 0, 0), card, camera)

		// Render the frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
