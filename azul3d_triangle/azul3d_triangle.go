// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays a few colored triangles.
package main

import (
	"fmt"
	"go/build"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"azul3d.org/gfx.v2-dev"
	"azul3d.org/gfx.v2-dev/window"
	"azul3d.org/keyboard.v1"
	math "azul3d.org/lmath.v1"
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

var (
	// Whether or not we should print the number of samples the triangle drew.
	printSamples bool

	// Whether or not we should print if the triangle's center is within the
	// camera's view.
	printInView bool
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, d gfx.Device) {
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camFOV := 75.0
	camNear := 0.0001
	camFar := 1000.0
	camera.SetPersp(d.Bounds(), camFOV, camNear, camFar)

	// Move the camera -2 on the Y axis (back two units away from the triangle
	// object).
	camera.SetPos(math.Vec3{0, -2, 0})

	// Loading shader files
	glslVert, err := ioutil.ReadFile(absPath("azul3d_triangle/triangle.vert"))
	if err != nil {
		log.Fatal(err)
	}
	glslFrag, err := ioutil.ReadFile(absPath("azul3d_triangle/triangle.frag"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSL = &gfx.GLSLSources{
		Vertex:   glslVert,
		Fragment: glslFrag,
	}

	// Create a triangle mesh.
	triMesh := gfx.NewMesh()
	triMesh.Vertices = []gfx.Vec3{
		// Top
		{0, 0, 1},
		{-.5, 0, 0},
		{.5, 0, 0},

		// Bottom-Left
		{-.5, 0, 0},
		{-1, 0, -1},
		{0, 0, -1},

		// Bottom-Right
		{.5, 0, 0},
		{0, 0, -1},
		{1, 0, -1},
	}
	triMesh.Colors = []gfx.Color{
		// Top
		{1, 0, 0, 1},
		{0, 1, 0, 1},
		{0, 0, 1, 1},

		// Bottom-Left
		{1, 0, 0, 1},
		{0, 1, 0, 1},
		{0, 0, 1, 1},

		// Bottom-Right
		{1, 0, 0, 1},
		{0, 1, 0, 1},
		{0, 0, 1, 1},
	}

	// Create a triangle object.
	triangle := gfx.NewObject()
	triangle.Shader = shader
	triangle.OcclusionTest = true
	triangle.State.FaceCulling = gfx.NoFaceCulling
	triangle.Meshes = []*gfx.Mesh{triMesh}

	// Transforms from different objects can be parented to one another to
	// create complex transformations (in this case we rotate -45 degrees then
	// +45 degrees which performs no rotation at all.
	right := gfx.NewTransform()
	right.SetRot(math.Vec3{0, 0, -45})

	left := gfx.NewTransform()
	left.SetRot(math.Vec3{0, 0, 45})
	left.SetParent(right)

	triangle.Transform.SetParent(left)

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
				camera.SetPersp(d.Bounds(), camFOV, camNear, camFar)
				camera.Unlock()

			case keyboard.TypedEvent:
				switch ev.Rune {
				case 's':
					printSamples = !printSamples
				case 'v':
					printInView = !printInView

				case 'm':
					// Toggle MSAA now.
					msaa := !d.MSAA()
					d.SetMSAA(msaa)
					fmt.Println("MSAA Enabled?", msaa)

				case 'p':
					triMesh.Lock()
					triMesh.Primitive = gfx.Points
					triMesh.Unlock()

				case 't':
					triMesh.Lock()
					triMesh.Primitive = gfx.Triangles
					triMesh.Unlock()

				case 'l':
					triMesh.Lock()
					triMesh.Primitive = gfx.Lines
					triMesh.Unlock()

				case '1':
					// Take a screenshot.
					fmt.Println("Writing screenshot to file...")
					// Download the image from the graphics hardware and save
					// it to disk.
					complete := make(chan image.Image, 1)
					d.Download(image.Rect(256, 256, 512, 512), complete)
					img := <-complete // Wait for download to complete.

					// Save to png.
					f, err := os.Create("screenshot.png")
					if err != nil {
						log.Fatal(err)
					}
					err = png.Encode(f, img)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Wrote texture to screenshot.png")

				case '2':
					fmt.Println("toggle fullscreen")
					props := w.Props()
					props.SetFullscreen(!props.Fullscreen())
					w.Request(props)
				}
			}
		}
	}()

	for {
		var v math.Vec2
		// Depending on keyboard state, transform the triangle.
		kb := w.Keyboard()
		if kb.Down(keyboard.ArrowLeft) {
			v.X -= 1
		}
		if kb.Down(keyboard.ArrowRight) {
			v.X += 1
		}
		if kb.Down(keyboard.ArrowDown) {
			v.Y -= 1
		}
		if kb.Down(keyboard.ArrowUp) {
			v.Y += 1
		}

		// Apply movement relative to the frame rate.
		v = v.MulScalar(d.Clock().Dt())

		// Update the triangle's transformation matrix.
		triangle.RLock()
		if kb.Down(keyboard.R) {
			// Reset transformation.
			oldParent := triangle.Transform.Parent()
			triangle.Transform.Reset()
			triangle.Transform.SetParent(oldParent)

		} else if kb.Down(keyboard.RightAlt) {
			// Apply shearing on X/Y axis.
			s := math.Vec3{v.Y, v.X, 0}
			if kb.Down(keyboard.RightShift) {
				// Apply shearing on X/Z axis.
				s = math.Vec3{v.Y, 0, v.X}
			}
			triangle.SetShear(triangle.Shear().Add(s.MulScalar(3)))

		} else if kb.Down(keyboard.LeftAlt) {
			// Apply scaling on X/Z axis.
			s := math.Vec3{v.X, 0, v.Y}
			if kb.Down(keyboard.LeftShift) {
				// Apply scaling on X/Y axis.
				s = math.Vec3{v.X, v.Y, 0}
			}
			triangle.SetScale(triangle.Scale().Add(s.MulScalar(3)))

		} else if kb.Down(keyboard.LeftCtrl) {
			// Apply rotation on X/Z axis.
			r := math.Vec3{v.Y, 0, v.X}
			if kb.Down(keyboard.LeftShift) {
				// Apply rotation on X/Y axis.
				r = math.Vec3{v.Y, v.X, 0}
			}
			triangle.SetRot(triangle.Rot().Add(r.MulScalar(90)))

		} else {
			// Apply movement on X/Z axis.
			p := math.Vec3{v.X, 0, v.Y}
			if kb.Down(keyboard.LeftShift) {
				// Apply movement on X/Y axis.
				p = math.Vec3{v.X, v.Y, 0}
			}
			triangle.SetPos(triangle.Pos().Add(p.MulScalar(3)))
		}
		triangle.RUnlock()

		// Clear color and depth buffers.
		d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
		d.ClearDepth(d.Bounds(), 1.0)

		// Clear a few rectangles on the window using different background
		// colors.
		d.Clear(image.Rect(0, 100, 720, 380), gfx.Color{0, 1, 0, 1})
		d.Clear(image.Rect(100, 100, 620, 380), gfx.Color{1, 0, 0, 1})
		d.Clear(image.Rect(100, 200, 620, 280), gfx.Color{0, 0.5, 0.5, 1})
		d.Clear(image.Rect(200, 200, 520, 280), gfx.Color{1, 1, 0, 1})

		// Draw the triangle to the screen.
		bounds := d.Bounds()
		d.Draw(bounds.Inset(50), triangle, camera)

		// Render the whole frame.
		d.Render()

		// Print the number of samples the triangle drew (only if the GPU
		// supports occlusion queries).
		if printSamples && d.Info().OcclusionQuery {
			// The number of samples the triangle drew:
			samples := triangle.SampleCount()

			// The number of pixels the triangle drew:
			msaa := d.Precision().Samples
			if msaa == 0 {
				msaa = 1
			}
			pixels := samples / msaa

			// The percent of the window that the triangle drew to:
			bounds := d.Bounds()
			percentage := float64(pixels) / float64(bounds.Dx()*bounds.Dy())

			fmt.Printf("Drew %v samples (%vpx, %f%% of window)\n", samples, pixels, percentage)
		}

		// Print if the triangle's center is in the view of the camera.
		if printInView {
			triangle.RLock()
			tp := triangle.Pos()
			triangle.RUnlock()

			camera.RLock()
			proj, ok := camera.Project(tp)
			fmt.Println("In view?", ok, proj)
			camera.RUnlock()
		}
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
