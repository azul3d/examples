// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays GPU information.
package main

import (
	"azul3d.org/chippy.v1"
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v1"
	"fmt"
)

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	defer w.Destroy()

	gpu := r.GPUInfo()
	fmt.Println("GPU Name:", gpu.Name)
	fmt.Println("GPU Vendor:", gpu.Vendor)
	fmt.Printf("OpenGL: v%d.%d\n", gpu.GLMajor, gpu.GLMinor)
	fmt.Printf("GLSL: v%d.%d\n", gpu.GLSLMajor, gpu.GLSLMinor)
	fmt.Println("OcclusionQuery =", gpu.OcclusionQuery)
	fmt.Println("OcclusionQueryBits =", gpu.OcclusionQueryBits)
	fmt.Println("MaxTextureSize =", gpu.MaxTextureSize)
	fmt.Println("NPOT Textures =", gpu.NPOT)
	fmt.Println("AlphaToCoverage =", gpu.AlphaToCoverage)
	fmt.Println("GLSLMaxVaryingFloats =", gpu.GLSLMaxVaryingFloats)
	fmt.Println("GLSLMaxVertexInputs =", gpu.GLSLMaxVertexInputs)
	fmt.Println("GLSLMaxFragmentInputs =", gpu.GLSLMaxFragmentInputs)
	fmt.Println("OpenGL extensions:", gpu.GLExtensions)
}

func main() {
	window.Run(gfxLoop)
}
