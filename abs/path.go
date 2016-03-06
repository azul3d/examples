// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package abs finds an absolute path to a example resource file.
package abs

import (
	"go/build"
	"os"
	"path/filepath"
	"sync"
)

var (
	pathLock    sync.Mutex
	examplesDir string
)

// Path returns the absolute path to a file given a relative one in the
// examples directory:
//
//  $GOPATH/src/azul3d.org/examples.v1
//
// This helper function is not an important concept to any of the examples, it
// just allows the examples to be ran from any working directory.
func Path(relPath string) string {
	pathLock.Lock()
	defer pathLock.Unlock()

	if len(examplesDir) == 0 {
		// Find assets directory.
		for _, path := range filepath.SplitList(build.Default.GOPATH) {
			path = filepath.Join(path, "src/azul3d.org/examples")
			if _, err := os.Stat(path); err == nil {
				examplesDir = path
				break
			}
		}
	}
	return filepath.Join(examplesDir, relPath)
}
