// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays a picture using the stencil buffer with shapes.
package main

import (
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/gfxutil"
	"azul3d.org/engine/gfx/window"
	math "azul3d.org/engine/lmath"

	"azul3d.org/examples/abs"
)

// cardMesh creates and returns a card mesh.
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

// cardTexCoords returns a slice of texture coordinates for a card given
// [u,v,s,t] texture coordinates.
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

// shapeTexCoords returns the correct set of texture coordinates for a shape
// index. The texture has four shapes (0, 1, 2, 3) so given an index we return
// texture coordinates selecting the correct region (i.e. top left, bottom
// right, etc) of the texture.
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
	default:
		panic("never here")
	}
}

var (
	shapeMeshCache = make(map[int]*gfx.Mesh)
	textureCache   = make(map[string]*gfx.Texture)
)

func loadShapeMesh(which int) *gfx.Mesh {
	// Check if that texture is already loaded.
	m, ok := shapeMeshCache[which]
	if ok {
		return m
	}

	// Create a card object.
	m = cardMesh(1.0, 1.0)
	m.TexCoords = []gfx.TexCoordSet{
		{
			Slice: shapeTexCoords(which),
		},
	}

	// Cache for later and return.
	shapeMeshCache[which] = m
	return m
}

// loadTexture loads the given image file and returns a texture. The result is
// cached such that if you try to load the same image file twice, you only get
// one texture pointer rather than two.
func loadTexture(path string) *gfx.Texture {
	// If the cache has the texture, return it.
	tex, ok := textureCache[path]
	if ok {
		return tex
	}

	// Cache does not have the texture, open it now and store it in the cache
	// for later.
	tex, err := gfxutil.OpenTexture(path)
	if err != nil {
		log.Fatal(err)
	}
	tex.Format = gfx.DXT1RGBA
	textureCache[path] = tex
	return tex
}

var shapeState = &gfx.State{
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

func createShape(d gfx.Device, path string, which int) *gfx.Object {
	// Create the object.
	card := gfx.NewObject()
	card.Textures = []*gfx.Texture{loadTexture(path)}

	// Load the shape's mesh.
	card.Meshes = []*gfx.Mesh{loadShapeMesh(which)}

	// Set the card's state.
	card.State = shapeState
	return card
}

var shapes []*gfx.Object

// Tells if the shape is within twice the window's size or not. We use twice
// the size to account for the largeness of the shape.
func isDead(cam *camera.Camera, shape *gfx.Object) bool {
	worldPos := shape.ConvertPos(shape.Pos(), gfx.LocalToWorld)
	viewPos, _ := cam.Project(worldPos)
	xValid := viewPos.X < 2 && viewPos.X > -2
	yValid := viewPos.Y < 2
	if !xValid || !yValid {
		return true
	}
	return false
}

var (
	butterfly = time.Tick(time.Second / 4)
	other     = time.Tick(time.Second / 2)
)

func updateShapes(d gfx.Device, shader *gfx.Shader, cam *camera.Camera) {
	// Butterfly.
	which := 0

	select {
	default:
		return
	case <-butterfly:
	case <-other:
		which = int(rand.Float64() * 4)
	}

	// Create a shape.
	shape := createShape(d, abs.Path("azul3d_stencil/shapes.png"), which)
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

	// Remove dead shapes.
	n := len(shapes)
	i := 0
l:
	for i < n {
		if isDead(cam, shapes[i]) {
			// Release object for re-use.
			shapes[i].Destroy()

			// Remove from slice.
			shapes[i] = shapes[n-1]
			n--
			continue l
		}
		i++
	}
	shapes = shapes[:n]
	shapes = append(shapes, shape)
}

// createBackground creates the background picture. The returned object has no
// shader assigned to it.
func createBackground() *gfx.Object {
	// Open the background texture.
	tex, err := gfxutil.OpenTexture(abs.Path("azul3d_stencil/yi_han_cheol.png"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a card object.
	aspect := float32(tex.Bounds.Dx()) / float32(tex.Bounds.Dy())
	var height float32 = 1.0
	cardMesh := cardMesh(aspect, height)
	cardMesh.TexCoords = []gfx.TexCoordSet{
		{
			Slice: cardTexCoords(0, 0, 1, 1),
		},
	}
	card := gfx.NewObject()
	card.Textures = []*gfx.Texture{tex}
	card.Meshes = []*gfx.Mesh{cardMesh}
	card.State = gfx.NewState()
	card.State.StencilTest = true
	card.State.StencilFront = gfx.StencilState{
		ReadMask:  0xFF,
		Reference: 1,
		Fail:      gfx.SZero,
		DepthFail: gfx.SZero,
		DepthPass: gfx.SKeep,
		Cmp:       gfx.Equal,
	}
	return card
}

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, d gfx.Device) {
	// This example requires a stencil buffer, if we didn't get one from the
	// device then we simply exit.
	if d.Precision().StencilBits == 0 {
		log.Fatal("Could not aquire a stencil buffer.")
	}

	// Create a new orthographic (2D) camera.
	cam := camera.NewOrtho(d.Bounds())

	// Move the camera back two units away from the card.
	cam.SetPos(math.Vec3{0, -2, 0})

	// updateCamera is a function used to update the camera's projection taking
	// into account the new window size.
	updateCamera := func() {
		bounds := d.Bounds()
		xRatio := float64(bounds.Dx()) / float64(bounds.Dy())
		m := math.Mat4Ortho(-xRatio, xRatio, -1, 1, cam.Near, cam.Far)
		cam.P = gfx.ConvertMat4(m)
	}
	updateCamera()

	// Read the GLSL shaders from disk.
	shader, err := gfxutil.OpenShader(abs.Path("azul3d_stencil/stencil"))
	if err != nil {
		log.Fatal(err)
	}

	// Create the background.
	bgPicture := createBackground()
	bgPicture.Shader = shader

	// We want to know when the framebuffer is resized so we can update our
	// camera.
	events := make(chan window.Event, 64)
	w.Notify(events, window.FramebufferResizedEvents)

	for {
		// Handle each pending FramebufferResized event.
		window.Poll(events, func(e window.Event) {
			// The framebuffer was resized, so we update our camera now.
			updateCamera()
		})

		updateShapes(d, shader, cam)

		// Clear the color, depth, and stencil buffers.
		d.Clear(d.Bounds(), gfx.Color{0, 0, 0, 1})
		d.ClearDepth(d.Bounds(), 1.0)
		d.ClearStencil(d.Bounds(), 0)

		for _, shape := range shapes {
			// Skip drawing of shapes that are dead.
			if isDead(cam, shape) {
				continue
			}

			// We will move the shape forward a small amount.
			v := math.Vec3{0, 0, 0.7 * d.Clock().Dt()}

			// We don't want movement to take scale into account, all shapes
			// move the same speed no matter how large or small.
			v = v.Mul(math.Vec3One.Div(shape.Scale()))

			// The vector, v, has the shape moving foward a bit relative to
			// it's current rotation/etc, we convert the vector from the
			// shape's local space (i.e. one-foot-forward) into world space (an
			// actual point).
			shape.SetPos(shape.ConvertPos(v, gfx.LocalToWorld))

			// Draw the shape.
			d.Draw(d.Bounds(), shape, cam)
		}

		// Draw the background picture.
		d.Draw(d.Bounds(), bgPicture, cam)

		// Render the whole frame.
		d.Render()
	}
}

func main() {
	props := window.NewProps()
	props.SetSize(720, 480)

	// Here we will explicitly request a 8-bit stencil buffer, we must have a
	// stencil buffer for this example to work properly!
	p := window.DefaultProps.Precision()
	p.StencilBits = 8
	props.SetPrecision(p)

	window.Run(gfxLoop, props)
}
