// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays GPU information.
package main

import (
	"fmt"

	"azul3d.org/gfx.v2-dev"
	"azul3d.org/gfx.v2-dev/window"
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, d gfx.Device) {
	defer w.Close()

	dev := d.Info()
	fmt.Println("Device Name:", dev.Name)
	fmt.Println("Device Vendor:", dev.Vendor)
	fmt.Printf("OpenGL: v%d.%d\n", dev.GLMajor, dev.GLMinor)
	fmt.Printf("GLSL: v%d.%d\n", dev.GLSLMajor, dev.GLSLMinor)
	fmt.Println("OcclusionQuery =", dev.OcclusionQuery)
	fmt.Println("OcclusionQueryBits =", dev.OcclusionQueryBits)
	fmt.Println("MaxTextureSize =", dev.MaxTextureSize)
	fmt.Println("NPOT Textures =", dev.NPOT)

	fmt.Printf("%d Render-To-Texture MSAA Formats:\n", len(dev.RTTFormats.Samples))
	for i, sampleCount := range dev.RTTFormats.Samples {
		fmt.Printf("    %d. %dx MSAA\n", i+1, sampleCount)
	}

	fmt.Printf("%d Render-To-Texture Color Formats:\n", len(dev.RTTFormats.ColorFormats))
	for i, f := range dev.RTTFormats.ColorFormats {
		fmt.Printf("    %d. %+v\n", i+1, f)
	}
	fmt.Printf("%d Render-To-Texture Depth Formats:\n", len(dev.RTTFormats.DepthFormats))
	for i, f := range dev.RTTFormats.DepthFormats {
		fmt.Printf("    %d. %+v\n", i+1, f)
	}
	fmt.Printf("%d Render-To-Texture Stencil Formats:\n", len(dev.RTTFormats.StencilFormats))
	for i, f := range dev.RTTFormats.StencilFormats {
		fmt.Printf("    %d. %+v\n", i+1, f)
	}

	fmt.Println("AlphaToCoverage =", dev.AlphaToCoverage)
	fmt.Println("GLSLMaxVaryingFloats =", dev.GLSLMaxVaryingFloats)
	fmt.Println("GLSLMaxVertexInputs =", dev.GLSLMaxVertexInputs)
	fmt.Println("GLSLMaxFragmentInputs =", dev.GLSLMaxFragmentInputs)
	fmt.Println("OpenGL extensions:", dev.GLExtensions)
}

func main() {
	props := window.NewProps()
	props.SetVisible(false)
	window.Run(gfxLoop, props)
}
