// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build qml

// Example - Uses QML to display graphics.
package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/gl2.v1"
	"gopkg.in/qml.v0"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

var (
	renderer *gl2.Renderer
	printFPS = time.Tick(1 * time.Second)
)

type GoRect struct {
	qml.Object
}

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the QML draw loop.
func gfxLoop(r gfx.Renderer) {
	for {
		select {
		case <-printFPS:
			cl := r.Clock()
			fmt.Printf("%v FPS (%v Avg.)\n", cl.FrameRate(), int(cl.AverageFrameRate()))
		default:
		}

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})

		// Clear a few rectangles on the GoRect using different background
		// colors.
		r.Clear(image.Rect(0, 10, 100, 90), gfx.Color{0, 1, 0, 1})
		r.Clear(image.Rect(10, 10, 90, 90), gfx.Color{1, 0, 0, 1})
		r.Clear(image.Rect(10, 20, 90, 80), gfx.Color{0, 0.5, 0.5, 1})
		r.Clear(image.Rect(20, 20, 80, 80), gfx.Color{1, 1, 0, 1})

		// Render the whole frame.
		r.Render()
	}
}

func (r *GoRect) Paint(p *qml.Painter) {
	if renderer == nil {
		// Initialize renderer.
		var err error
		renderer, err = gl2.New(true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Write renderer debug output (shader errors, etc) to stdout.
		renderer.SetDebugOutput(os.Stdout)

		// Start the graphics rendering loop.
		go gfxLoop(renderer)
	}

	// The QML rectangle's size may have changed (i.e. animated size), so
	// update the renderer's bounds.
	renderer.UpdateBounds(image.Rect(0, 0, r.Int("width"), r.Int("height")))

	// Execute OpenGL commands until a frame has been rendered.
	for {
		select {
		case fn := <-renderer.RenderExec:
			if renderedFrame := fn(); renderedFrame {
				// Request that QML give us another frame later.
				r.Call("update")
				return
			}

		case fn := <-renderer.LoaderExec:
			fn()
		}
	}
}

func run() error {
	qml.Init(nil)

	qml.RegisterTypes("GoExtensions", 1, 0, []qml.TypeSpec{{
		Init: func(r *GoRect, obj qml.Object) { r.Object = obj },
	}})

	engine := qml.NewEngine()
	component, err := engine.LoadFile("src/azul3d.org/v1/examples/azul3d_qml_clearing/azul3d_qml_clearing.qml")
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)
	win.Show()
	win.Wait()

	return nil
}
