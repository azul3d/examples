// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Uses termbox to visualize D* Lite grid pathfinding.
package main

import (
	"log"

	"azul3d.org/engine/dstarlite/grid"
	"github.com/nsf/termbox-go"
)

type Player struct {
	X, Y int
}

var (
	player Player
	g      *grid.Data
)

func draw() {
	width, height := g.Size()

	// Draw grid as seen by DSL
	for x := 0; x < width*2; x++ {
		for y := 0; y < height; y++ {
			for n := 0; n < 2; n++ {
				color := termbox.ColorWhite
				v, ok := g.Get(grid.Coord{x / 2, y})
				if !ok {
					log.Fatal("This shouldn't happen (coordinate outside grid).")
				}

				if v == -1 {
					color = termbox.ColorRed
				}
				termbox.SetCell(x+n, y, rune(' '), termbox.ColorDefault, color)
			}
		}
	}

	// Draw path
	for _, coord := range g.Plan() {
		for n := 0; n < 2; n++ {
			termbox.SetCell(coord[0]*2+n, coord[1], rune(' '), termbox.ColorDefault, termbox.ColorYellow)
		}
	}

	// Draw player and goal
	for n := 0; n < 2; n++ {
		termbox.SetCell(player.X*2+n, player.Y, rune(' '), termbox.ColorDefault, termbox.ColorBlue)

		goal := g.Goal()
		termbox.SetCell(goal[0]*2+n, goal[1], rune(' '), termbox.ColorDefault, termbox.ColorGreen)
	}

	// Draw path
	termbox.Flush()
}

func main() {
	// Init termbox
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	// We use half the terminal width (since approx 2 characters wide is 'square')
	width, height := termbox.Size()
	width /= 2

	// Create grid
	g = grid.New(width, height, grid.Coord{0, 0}, grid.Coord{width / 2, height / 2})

	// Resets the player position and clears all grid cells
	reset := func() {
		player.X = 0
		player.Y = 0
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				g.Set(grid.Coord{x, y}, 1)
			}
		}
		draw()
	}
	reset()

	// Main loop to wait for keyboard input
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyArrowUp:
				if player.Y > 0 {
					player.Y -= 1
				}
				draw()

			case ev.Key == termbox.KeyArrowDown:
				if player.Y < height-1 {
					player.Y += 1
				}
				draw()

			case ev.Key == termbox.KeyArrowLeft:
				if player.X > 0 {
					player.X -= 1
				}
				draw()

			case ev.Key == termbox.KeyArrowRight:
				if player.X < width-1 {
					player.X += 1
				}
				draw()

			case ev.Key == termbox.KeyEsc:
				break loop

			case ev.Key == termbox.KeyEnter:
				draw()

			case ev.Key == termbox.KeySpace:
				c := grid.Coord{player.X, player.Y}
				v, ok := g.Get(c)
				if ok {
					if v == -1 {
						g.Set(c, 1)
					} else {
						g.Set(c, -1)
					}
				}
				draw()

			case ev.Ch == rune('s') || ev.Ch == rune('S'):
				g.UpdateStart(grid.Coord{player.X, player.Y})
				draw()

			case ev.Ch == rune('r') || ev.Ch == rune('R'):
				reset()
				draw()
			}
		}
	}
}
