// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build glfw

// Example - Uses GLFW to display graphics.
package main

import (
	"fmt"
	"image"
	"os"
	"runtime"
	"time"

	"azul3d.org/clock.v1"
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/gl2.v2"
	"github.com/go-gl/glfw3"
)

func errorCallback(err glfw3.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the GLFW main loop.
func gfxLoop(r gfx.Renderer) {
	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})

		// Clear a few rectangles on the window using different background
		// colors.
		r.Clear(image.Rect(0, 100, 640, 380), gfx.Color{0, 1, 0, 1})
		r.Clear(image.Rect(100, 100, 540, 380), gfx.Color{1, 0, 0, 1})
		r.Clear(image.Rect(100, 200, 540, 280), gfx.Color{0, 0.5, 0.5, 1})
		r.Clear(image.Rect(200, 200, 440, 280), gfx.Color{1, 1, 0, 1})

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	glfw3.SetErrorCallback(errorCallback)

	if !glfw3.Init() {
		panic("Can't init glfw!")
	}
	defer glfw3.Terminate()

	window, err := glfw3.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	// Make the context current.
	window.MakeContextCurrent()
	defer glfw3.DetachCurrentContext()

	// Disable vertical sync.
	glfw3.SwapInterval(0)

	// Create the OpenGL 2.0 renderer. A error is returned if the OpenGL
	// version is invalid.
	r, err := gl2.New()
	if err != nil {
		panic(err)
	}

	// Write renderer debug output (shader errors, etc) to stdout.
	r.SetDebugOutput(os.Stdout)

	// Whenever the window is resized, inform the renderer that it's bounds
	// have changed.
	window.SetSizeCallback(func(w *glfw3.Window, width, height int) {
		r.UpdateBounds(image.Rect(0, 0, width, height))
	})

	// Start the graphics rendering loop.
	go gfxLoop(r)

	cl := clock.New()
	cl.SetMaxFrameRate(0)
	printFPS := time.Tick(1 * time.Second)

	for !window.ShouldClose() {
		select {
		case <-printFPS:
			fmt.Printf("%v FPS (%v Avg.)\n", cl.FrameRate(), int(cl.AverageFrameRate()))

		case fn := <-r.RenderExec:
			if renderedFrame := fn(); renderedFrame {
				// Tell the clock a new frame has begun.
				cl.Tick()

				// Swap OpenGL buffers.
				window.SwapBuffers()

				// Poll for GLFW events.
				glfw3.PollEvents()
			}

		case fn := <-r.LoaderExec:
			fn()
		}
	}
}
