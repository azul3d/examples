// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates debug camera. Use "t" to toggle camera modes.
package main

import (
	"image"
	"log"
	"math"

	"azul3d.org/engine/gfx"
	"azul3d.org/engine/gfx/camera"
	"azul3d.org/engine/gfx/gfxutil"
	"azul3d.org/engine/gfx/window"
	"azul3d.org/engine/keyboard"
	"azul3d.org/engine/lmath"

	"azul3d.org/examples/abs"
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
func setup(d gfx.Device) (camMain, camSecondary *camera.Camera, cubes *gfx.Object) {
	// Load the shader.
	shader, err := gfxutil.OpenShader(abs.Path("azul3d_debug_camera/cube"))
	if err != nil {
		log.Fatal(err)
	}

	// Create and position both cameras.
	camMain = camera.New(d.Bounds())
	camMain.Debug = true // Enable drawing camera as wireframe.
	camSecondary = camera.New(d.Bounds())
	camSecondary.SetPos(lmath.Vec3{10, -10, 10})
	camSecondary.SetRot(lmath.Vec3{-45, 0, 45})

	// Insert some cubes so we have something to look at.
	cubes = gfx.NewObject()
	cubes.State = gfx.NewState()
	cubes.State.FaceCulling = gfx.BackFaceCulling
	cubes.Shader = shader
	cubes.Meshes = []*gfx.Mesh{}
	for x := -3; x <= 3; x++ {
		for y := -3; y <= 3; y++ {
			cubes.Meshes = append(cubes.Meshes, cube(float32(x)*1.5, float32(y)*1.5))
		}
	}
	return
}

func update(d gfx.Device, cubes *gfx.Object, camMain *camera.Camera) {
	// Clear the entire area.
	d.Clear(d.Bounds(), gfx.Color{0, 0, 0, 1})
	d.ClearDepth(d.Bounds(), 1.0)

	// Use frame delta time so the windows stay in sync.
	dt := d.Clock().Dt()
	cubes.SetRot(cubes.Rot().AddScalar(10 * dt))

	// How much we want to add to the FOV and far clipping distance (pulsating).
	add := math.Sin(float64(d.Clock().FrameCount())/20) * 25

	camMain.Near = 1
	camMain.Far = 30 + add

	if camMain.Ortho {
		// Reduce the viewing area so the cubes don't get too small.
		x := d.Bounds().Dx() / 48
		y := d.Bounds().Dy() / 48

		// Set camera position and orthographic projection settings.
		camMain.Update(image.Rect(0, 0, x, y))
		camMain.SetPos(lmath.Vec3{-float64(x) / 2, -10, -float64(y) / 2})
		return
	}

	// Set camera perspective and update the position of the camera in case
	// we were previously using an orthographic camera.
	camMain.FOV = 75 + add
	camMain.Update(d.Bounds())
	camMain.SetPos(lmath.Vec3{0, -10, 0})
}

// Used by one window(2) to signal to the other(1) the current camera mode.
var camOrtho = make(chan bool, 1)

// gfxLoopWindow1 runs the main camera window
func gfxLoopWindow1(w window.Window, d gfx.Device) {
	camMain, _, cubes := setup(d)
	for {
		select {
		case ortho := <-camOrtho:
			camMain.Ortho = ortho
		default:
		}

		// Update the camera position, rotation etc (done in both windows).
		update(d, cubes, camMain)
		d.Draw(d.Bounds(), cubes, camMain)

		// Render the whole frame.
		d.Render()
	}
}

// gfxLoopWindow2 runs the second camera, observing the camera helper
func gfxLoopWindow2(w window.Window, d gfx.Device) {
	// Have the window notify us whenever typing events occur.
	events := make(chan window.Event, 256)
	w.Notify(events, window.KeyboardTypedEvents)

	camMain, camSecondary, cubes := setup(d)
	for {
		select {
		case e := <-events:
			switch ev := e.(type) {
			case keyboard.Typed:
				switch ev.S {
				case "t":
					camMain.Ortho = !camMain.Ortho
					camOrtho <- camMain.Ortho
				}
			}
		default:
		}

		// Update the camera position, rotation etc (done in both windows).
		update(d, cubes, camMain)

		// Draw the first camera as seen by the secondary camera.
		d.Draw(d.Bounds(), camMain.Object, camSecondary)

		// Draw the cubes as seen by the secondary camera.
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
