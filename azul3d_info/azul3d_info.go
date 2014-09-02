// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays GPU information.
package main

import (
	"fmt"

	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, r gfx.Renderer) {
	defer w.Close()

	gpu := r.GPUInfo()
	fmt.Println("GPU Name:", gpu.Name)
	fmt.Println("GPU Vendor:", gpu.Vendor)
	fmt.Printf("OpenGL: v%d.%d\n", gpu.GLMajor, gpu.GLMinor)
	fmt.Printf("GLSL: v%d.%d\n", gpu.GLSLMajor, gpu.GLSLMinor)
	fmt.Println("OcclusionQuery =", gpu.OcclusionQuery)
	fmt.Println("OcclusionQueryBits =", gpu.OcclusionQueryBits)
	fmt.Println("MaxTextureSize =", gpu.MaxTextureSize)
	fmt.Println("NPOT Textures =", gpu.NPOT)

	fmt.Printf("%d Render-To-Texture MSAA Formats:\n", len(gpu.RTTFormats.Samples))
	for i, sampleCount := range gpu.RTTFormats.Samples {
		fmt.Printf("    %d. %dx MSAA\n", i+1, sampleCount)
	}

	fmt.Printf("%d Render-To-Texture Color Formats:\n", len(gpu.RTTFormats.ColorFormats))
	for i, f := range gpu.RTTFormats.ColorFormats {
		fmt.Printf("    %d. %+v\n", i+1, f)
	}
	fmt.Printf("%d Render-To-Texture Depth Formats:\n", len(gpu.RTTFormats.DepthFormats))
	for i, f := range gpu.RTTFormats.DepthFormats {
		fmt.Printf("    %d. %+v\n", i+1, f)
	}
	fmt.Printf("%d Render-To-Texture Stencil Formats:\n", len(gpu.RTTFormats.StencilFormats))
	for i, f := range gpu.RTTFormats.StencilFormats {
		fmt.Printf("    %d. %+v\n", i+1, f)
	}

	fmt.Println("AlphaToCoverage =", gpu.AlphaToCoverage)
	fmt.Println("GLSLMaxVaryingFloats =", gpu.GLSLMaxVaryingFloats)
	fmt.Println("GLSLMaxVertexInputs =", gpu.GLSLMaxVertexInputs)
	fmt.Println("GLSLMaxFragmentInputs =", gpu.GLSLMaxFragmentInputs)
	fmt.Println("OpenGL extensions:", gpu.GLExtensions)
}

func main() {
	props := window.NewProps()
	props.SetVisible(false)
	window.Run(gfxLoop, props)
}
