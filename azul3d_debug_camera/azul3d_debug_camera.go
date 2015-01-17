// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates debug camera. Use "t" to toggle camera modes.
package main

import (
	"image"
	"log"
	"math"

	"azul3d.org/gfx.v2-dev"
	"azul3d.org/gfx.v2-dev/debug"
	"azul3d.org/gfx.v2-dev/gfxutil"
	"azul3d.org/gfx.v2-dev/window"
	"azul3d.org/keyboard.v2-dev"
	"azul3d.org/lmath.v1"
)

var (
	perspCamera = true
)

// cube returns a cube *gfx.Mesh at an offset from the origin
func cube(x, y float32) *gfx.Mesh {
	m := gfx.NewMesh()
	m.Vertices = []gfx.Vec3{
		{-0.5 + x, -0.5 + y, 0.5},
		{0.5 + x, -0.5 + y, 0.5},
		{0.5 + x, 0.5 + y, 0.5},
		{-0.5 + x, 0.5 + y, 0.5},

		{-0.5 + x, -0.5 + y, -0.5},
		{-0.5 + x, 0.5 + y, -0.5},
		{0.5 + x, 0.5 + y, -0.5},
		{0.5 + x, -0.5 + y, -0.5},

		{-0.5 + x, 0.5 + y, -0.5},
		{-0.5 + x, 0.5 + y, 0.5},
		{0.5 + x, 0.5 + y, 0.5},
		{0.5 + x, 0.5 + y, -0.5},

		{-0.5 + x, -0.5 + y, -0.5},
		{0.5 + x, -0.5 + y, -0.5},
		{0.5 + x, -0.5 + y, 0.5},
		{-0.5 + x, -0.5 + y, 0.5},

		{0.5 + x, -0.5 + y, -0.5},
		{0.5 + x, 0.5 + y, -0.5},
		{0.5 + x, 0.5 + y, 0.5},
		{0.5 + x, -0.5 + y, 0.5},

		{-0.5 + x, -0.5 + y, -0.5},
		{-0.5 + x, -0.5 + y, 0.5},
		{-0.5 + x, 0.5 + y, 0.5},
		{-0.5 + x, 0.5 + y, -0.5},
	}

	m.Indices = []uint32{
		0, 1, 2, 0, 2, 3, // front
		4, 5, 6, 4, 6, 7, // back
		8, 9, 10, 8, 10, 11, // top
		12, 13, 14, 12, 14, 15, // bottom
		16, 17, 18, 16, 18, 19, // right
		20, 21, 22, 20, 22, 23, // left
	}

	return m
}

// setup creates the cameras and cubes for both windows
func setup(d gfx.Device) (*gfx.Camera, *gfx.Camera, *gfx.Object) {
	// Load the shader.
	shader, err := gfxutil.OpenShader("cube")
	if err != nil {
		log.Fatal(err)
	}

	// Create and position both cameras.
	camMain := gfx.NewCamera()
	camSecondary := gfx.NewCamera()
	camSecondary.SetPersp(d.Bounds(), 75, 0.1, 1000)
	camSecondary.SetPos(lmath.Vec3{10, -10, 10})
	camSecondary.SetRot(lmath.Vec3{-45, 0, 45})

	// Insert some cubes so we have something to look at.
	cubes := gfx.NewObject()
	cubes.State = gfx.NewState()
	cubes.State.FaceCulling = gfx.BackFaceCulling
	cubes.Shader = shader
	cubes.Meshes = []*gfx.Mesh{}
	for x := -3; x <= 3; x++ {
		for y := -3; y <= 3; y++ {
			cubes.Meshes = append(cubes.Meshes, cube(float32(x)*1.5, float32(y)*1.5))
		}
	}

	return camMain, camSecondary, cubes
}

func update(d gfx.Device, cubes *gfx.Object, camMain *gfx.Camera) {
	// Clear the entire area.
	d.Clear(d.Bounds(), gfx.Color{0, 0, 0, 1})
	d.ClearDepth(d.Bounds(), 1.0)

	// Use frame delta time so the windows stay in sync.
	dt := d.Clock().Dt()
	cubes.SetRot(cubes.Rot().AddScalar(10 * dt))

	// How much we want to add to the FOV and far clipping distance (pulsating).
	add := math.Sin(float64(d.Clock().FrameCount())/20) * 25

	if perspCamera {
		// Set camera perspective and update the position of the camera in case
		// we were previously using an orthographic camera.
		camMain.SetPersp(d.Bounds(), 75+add, 1, 30+add)
		camMain.SetPos(lmath.Vec3{0, -10, 0})
	} else {
		// Reduce the viewing area so the cubes don't get too small.
		x := d.Bounds().Dx() / 48
		y := d.Bounds().Dy() / 48

		// Set camera position and orthographic projection settings.
		camMain.SetOrtho(image.Rect(0, 0, x, y), 1, 30+add)
		camMain.SetPos(lmath.Vec3{-float64(x) / 2, -10, -float64(y) / 2})
	}
}

// gfxLoopWindow1 runs the main camera window
func gfxLoopWindow1(w window.Window, d gfx.Device) {
	camMain, _, cubes := setup(d)
	for {
		// Update the camera position, rotation etc (done in both windows).
		update(d, cubes, camMain)
		d.Draw(d.Bounds(), cubes, camMain)

		// Render the whole frame.
		d.Render()
	}
}

// gfxLoopWindow2 runs the second camera, observing the camera helper
func gfxLoopWindow2(w window.Window, d gfx.Device) {
	// Create an event mask for the events we are interested in.
	evMask := window.KeyboardTypedEvents

	// Create a channel of events.
	events := make(chan window.Event, 256)

	// Have the window notify our channel whenever events occur.
	w.Notify(events, evMask)

	camMain, camSecondary, cubes := setup(d)
	for {
		// Handle each pending event.
		window.Poll(events, func(e window.Event) {
			switch ev := e.(type) {
			case keyboard.Typed:
				switch ev.S {
				case "t":
					perspCamera = !perspCamera
				}
			}
		})

		// Update the camera position, rotation etc (done in both windows).
		update(d, cubes, camMain)

		debug.DrawCamera(d, camMain, camSecondary)
		d.Draw(d.Bounds(), cubes, camSecondary)

		// Render the whole frame.
		d.Render()

	}
}

func main() {
	go func() {
		// Create our windows.
		props := window.NewProps()
		props.SetTitle("Main camera {FPS}")
		props.SetPos(0, 0)
		props.SetSize(640, 400)
		w, r, err := window.New(props)
		if err != nil {
			log.Fatal(err)
		}

		props = window.NewProps()
		props.SetTitle("Observer {FPS}")
		props.SetPos(640, 0)
		props.SetSize(640, 400)
		w2, r2, err := window.New(props)
		if err != nil {
			log.Fatal(err)
		}

		go gfxLoopWindow1(w, r)
		go gfxLoopWindow2(w2, r2)
	}()

	// Enter the main loop.
	window.MainLoop()
}
