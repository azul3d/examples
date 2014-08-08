// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays a picture using the stencil buffer with shapes.
package main

import (
	"azul3d.org/chippy.v1"
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/gl2.v1"
	"azul3d.org/gfx/window.v1"
	math "azul3d.org/lmath.v1"
	"image"
	_ "image/jpeg"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
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
		for _, path := range filepath.SplitList(os.ExpandEnv("$GOPATH")) {
			path = filepath.Join(path, "src/azul3d.org/examples.v1")
			if _, err := os.Stat(path); err == nil {
				examplesDir = path
				break
			}
		}
	}
	return filepath.Join(examplesDir, relPath)
}

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

// Creates and returns a card mesh.
func cardMesh(w, h float32) *gfx.Mesh {
	m := gfx.NewMesh()
	m.Vertices = []gfx.Vec3{
		// Left triangle.
		{-w, 0, h},  // Left-Top
		{-w, 0, -h}, // Left-Bottom
		{w, 0, -h},  // Right-Bottom

		// Right triangle.
		{-w, 0, h}, // Left-Top
		{w, 0, -h}, // Right-Bottom
		{w, 0, h},  // Right-Top
	}
	return m
}

// Returns a slice of texture coordinates for a card given u,v,s,t coordinates.
func cardTexCoords(u, v, s, t float32) []gfx.TexCoord {
	return []gfx.TexCoord{
		// Left triangle.
		{u, v},
		{u, t},
		{s, t},

		// Right triangle.
		{u, v},
		{s, t},
		{s, v},
	}
}

func createPicture(r gfx.Renderer, path string) *gfx.Object {
	// Load the picture.
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Create new texture and ask the renderer to load it.
	tex := gfx.NewTexture()
	tex.Source = img
	tex.MinFilter = gfx.LinearMipmapLinear
	tex.MagFilter = gfx.Linear
	tex.Format = gfx.DXT1
	aspect := float32(img.Bounds().Dx()) / float32(img.Bounds().Dy())
	var height float32 = 1.0

	// Create a card object.
	cardMesh := cardMesh(aspect, height)
	cardMesh.TexCoords = []gfx.TexCoordSet{
		{
			Slice: cardTexCoords(0, 0, 1, 1),
		},
	}
	card := gfx.NewObject()
	card.Textures = []*gfx.Texture{tex}
	card.Meshes = []*gfx.Mesh{cardMesh}
	return card
}

func shapeTexCoords(index int) []gfx.TexCoord {
	switch index {
	case 0:
		return cardTexCoords(0, 0, .5, .5)
	case 1:
		return cardTexCoords(.5, 0, 1, 0)
	case 2:
		return cardTexCoords(0, .5, .5, 1)
	case 3:
		return cardTexCoords(.5, .5, 1, 1)
	}
	panic("never here")
}

func createShape(r gfx.Renderer, path string, which int) *gfx.Object {
	// Load the picture.
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Create new texture and ask the renderer to load it.
	tex := gfx.NewTexture()
	tex.Source = img
	tex.MinFilter = gfx.LinearMipmapLinear
	tex.MagFilter = gfx.Linear
	tex.Format = gfx.DXT1RGBA

	// Create a card object.
	cardMesh := cardMesh(1.0, 1.0)
	cardMesh.TexCoords = []gfx.TexCoordSet{
		{
			Slice: shapeTexCoords(which),
		},
	}
	card := gfx.NewObject()
	card.Textures = []*gfx.Texture{tex}
	card.Meshes = []*gfx.Mesh{cardMesh}

	// Set the card's state.
	card.State = gfx.State{
		AlphaMode:   gfx.AlphaToCoverage,
		WriteRed:    false,
		WriteGreen:  false,
		WriteBlue:   false,
		WriteAlpha:  false,
		DepthWrite:  false,
		StencilTest: true,
		StencilFront: gfx.StencilState{
			WriteMask: 0xFF,
			Fail:      gfx.SReplace,
			DepthFail: gfx.SReplace,
			DepthPass: gfx.SReplace,
			Cmp:       gfx.Always,
			Reference: 1,
		},
	}
	return card
}

var shapes struct {
	sync.Mutex
	slice []*gfx.Object
}

// Tells if the shape is within twice the window's size or not. We use twice
// the size to account for the largeness of the shape.
func isDead(camera *gfx.Camera, shape *gfx.Object) bool {
	worldPos := shape.ConvertPos(shape.Pos(), gfx.LocalToWorld)
	camera.RLock()
	viewPos, _ := camera.Project(worldPos)
	camera.RUnlock()
	xValid := viewPos.X < 2 && viewPos.X > -2
	yValid := viewPos.Y < 2
	if !xValid || !yValid {
		return true
	}
	return false
}

func shapeSpawner(r gfx.Renderer, shader *gfx.Shader, camera *gfx.Camera) {
	butterfly := time.Tick(time.Second / 4)
	other := time.Tick(time.Second / 2)

	for {
		// Butterfly.
		which := 0

		select {
		case <-butterfly:
		case <-other:
			which = int(rand.Float64() * 4)
		}

		// Create a shape.
		shape := createShape(r, absPath("assets/textures/shapes.png"), which)
		shape.Shader = shader
		shape.SetPos(math.Vec3{0, -1, 0})

		// Give the shape a random scale.
		var s float64
		if which == 0 {
			s = rand.Float64() * 0.42
			if s < 0.2 {
				s = 0.2
			}
		} else {
			// Stars and other things are smaller.
			s = rand.Float64() * 0.21
			if s < 0.1 {
				s = 0.1
			}
		}
		shape.SetScale(math.Vec3{s, s, s})

		// Give the shape a random position.
		x := (rand.Float64() * 2.0) - 1.0
		y := (rand.Float64() * 2.0) - 4.0
		shape.SetPos(math.Vec3{x, 0, y})

		// Give the shape a random rotation.
		r := ((rand.Float64() * 2.0) - 1.0) * 45
		shape.SetRot(math.Vec3{0, r, 0})

		shapes.Lock()

		// Remove dead shapes.
		n := len(shapes.slice)
		i := 0
	l:
		for i < n {
			if isDead(camera, shapes.slice[i]) {
				// Remove from slice.
				shapes.slice[i] = shapes.slice[n-1]
				n--
				continue l
			}
			i++
		}
		shapes.slice = shapes.slice[:n]

		shapes.slice = append(shapes.slice, shape)
		shapes.Unlock()
	}
}

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	w.SetSize(720, 480)
	w.SetPositionCenter(chippy.DefaultScreen())
	glr := r.(*gl2.Renderer)
	glr.UpdateBounds(image.Rect(0, 0, 720, 480))

	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSLVert = glslVert
	shader.GLSLFrag = glslFrag

	// Create the background.
	bgPicture := createPicture(r, absPath("assets/textures/yi_han_cheol.jpg"))
	bgPicture.Shader = shader

	// Create a camera.
	c := gfx.NewCamera()
	c.SetPos(math.Vec3{0, -2, 0})

	// Start the shape spawner.
	go shapeSpawner(r, shader, c)

	for {
		bounds := r.Bounds()
		xRatio := float64(bounds.Dx()) / float64(bounds.Dy())
		m := math.Mat4Ortho(-xRatio, xRatio, -1, 1, 0.001, 100.0)
		c.Lock()
		c.Projection = gfx.ConvertMat4(m)
		c.Unlock()

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{0, 0, 0, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)
		r.ClearStencil(image.Rect(0, 0, 0, 0), 0)

		shapes.Lock()
		for _, shape := range shapes.slice {
			// Skip drawing of shapes that are dead.
			if isDead(c, shape) {
				continue
			}

			// We will move the shape forward a small amount.
			v := math.Vec3{0, 0, 0.7 * r.Clock().Dt()}

			// We don't want movement to take scale into account, all shapes
			// move the same speed no matter how large or small.
			v = v.Mul(math.Vec3One.Div(shape.Scale()))

			// Convert the position to world space.
			shape.SetPos(shape.ConvertPos(v, gfx.LocalToWorld))

			// Draw the shape.
			r.Draw(image.Rect(0, 0, 0, 0), shape, c)
		}
		shapes.Unlock()

		// Draw the background picture.
		bgPicture.State.StencilTest = true
		bgPicture.State.StencilFront = gfx.StencilState{
			ReadMask:  0xFF,
			Reference: 1,
			Fail:      gfx.SZero,
			DepthFail: gfx.SZero,
			DepthPass: gfx.SKeep,
			Cmp:       gfx.Equal,
		}
		r.Draw(image.Rect(0, 0, 0, 0), bgPicture, c)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop)
}
