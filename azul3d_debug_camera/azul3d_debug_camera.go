// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates debug camera.
package main

import (
	"log"

	"azul3d.org/gfx.v2-dev"
	"azul3d.org/gfx.v2-dev/debug"
	"azul3d.org/gfx.v2-dev/gfxutil"
	"azul3d.org/gfx.v2-dev/window"
	"azul3d.org/lmath.v1"
)

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

func setup(d gfx.Device) (*gfx.Camera, *gfx.Camera, *gfx.Object) {
	shader, err := gfxutil.OpenShader("cube")
	if err != nil {
		log.Fatal(err)
	}

	camMain := gfx.NewCamera()
	camMain.SetPersp(d.Bounds(), 75, 0.1, 100)
	camMain.SetPos(lmath.Vec3{0, -5, 0})

	camSecondary := gfx.NewCamera()
	camSecondary.SetPersp(d.Bounds(), 75, 0.1, 100)
	camSecondary.SetPos(lmath.Vec3{10, -5, 10})
	camSecondary.SetRot(lmath.Vec3{45, 0, 45})

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

func gfxLoopWindow1(w window.Window, d gfx.Device) {
	camMain, _, cubes := setup(d)
	for {
		// Clear the entire area.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})

		cubes.SetRot(cubes.Rot().AddScalar(0.2))

		d.Draw(d.Bounds(), cubes, camMain)

		// Render the whole frame.
		d.Render()

	}
}

func gfxLoopWindow2(w window.Window, d gfx.Device) {
	camMain, camSecondary, cubes := setup(d)
	for {
		// Clear the entire area.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})

		cubes.SetRot(cubes.Rot().AddScalar(0.2))

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
		props.SetSize(400, 300)
		w, r, err := window.New(props)
		if err != nil {
			log.Fatal(err)
		}

		props = window.NewProps()
		props.SetTitle("Observer {FPS}")
		props.SetPos(400, 0)
		props.SetSize(400, 300)
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
