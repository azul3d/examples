// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates receiving window events.
package main

import (
	"fmt"
	"image"
	"reflect"

	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/keyboard.v1"
	"azul3d.org/mouse.v1"
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, r gfx.Renderer) {
	// You can handle window events in a seperate goroutine!
	go func() {
		// Create our events channel with sufficient buffer size.
		events := make(chan window.Event, 256)

		// Notify our channel anytime any event occurs.
		w.Notify(events, window.AllEvents)

		// Wait for events.
		for event := range events{
			// Use reflection to print the type of event:
			fmt.Println("Event type:", reflect.TypeOf(event))

			// Print the event:
			fmt.Println(event)
		}
	}()

	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})

		// The keyboard is monitored for you, simply check if a key is down:
		if w.Keyboard().Down(keyboard.Space) {
			// Clear a red rectangle.
			r.Clear(image.Rect(0, 0, 100, 100), gfx.Color{1, 0, 0, 1})
		}

		// And the same thing with the mouse, check if a mouse button is down:
		if w.Mouse().Down(mouse.Left) {
			// Clear a blue rectangle.
			r.Clear(image.Rect(100, 100, 200, 200), gfx.Color{0, 0, 1, 1})
		}

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
