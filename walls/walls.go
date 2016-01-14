// Copyright 2015-2016 Richard Hawkins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package walls TODO doc

package walls

import (
	"fmt"
	"runtime"

	"github.com/hurricanerix/transylvania/shapes"
	"github.com/hurricanerix/transylvania/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Wall struct {
	Image *sprite.Context
	Rect  *shapes.Rect
	dx    float32
}

// New TODO doc
func New(group *sprite.Group) (*Wall, error) {
	// TODO should take a group in as a argument
	w := Wall{}

	wall, err := sprite.Load("resistor.png", 32)
	if err != nil {
		return &w, fmt.Errorf("could not load wall: %v", err)
	}
	w.Image = wall

	rect, err := shapes.NewRect(640.0, 440.0, float32(w.Image.Width), float32(w.Image.Height))
	if err != nil {
		return &w, fmt.Errorf("could create rect: %v", err)
	}
	w.Rect = rect

	// TODO: this should probably be added outside of player
	group.Add(&w)
	return &w, nil
}

// Bind TODO doc
func (w *Wall) Bind(program uint32) error {
	return w.Image.Bind(program)
}

// Update TODO doc
func (w *Wall) Update(dt float32, g *sprite.Group) {
	w.dx = -250.0
	w.Rect.X += w.dx * dt
}

// Draw TODO doc
func (w *Wall) Draw() {
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y+92.0)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y+92.0)
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y+64.0)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y+64.0)
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y+32.0)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y+32.0)

	// Wires ^^^

	w.Image.DrawFrame(0, w.Rect.X, w.Rect.Y)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y)
	w.Image.DrawFrame(1, w.Rect.X, w.Rect.Y-32.0)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0)
	w.Image.DrawFrame(2, w.Rect.X, w.Rect.Y-32.0*2)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*2)
	w.Image.DrawFrame(3, w.Rect.X, w.Rect.Y-32.0*3)
	w.Image.DrawFrame(4, w.Rect.X+32.0, w.Rect.Y-32*3.0)

	// Hole

	// bottom resistor
	w.Image.DrawFrame(5, w.Rect.X, w.Rect.Y-32.0*6)
	w.Image.DrawFrame(6, w.Rect.X+32.0, w.Rect.Y-32.0*6)
	w.Image.DrawFrame(0, w.Rect.X, w.Rect.Y-32*7)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32*7)
	w.Image.DrawFrame(1, w.Rect.X, w.Rect.Y-32.0*8)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*8)
	w.Image.DrawFrame(2, w.Rect.X, w.Rect.Y-32.0*9)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*9)

	// bottom wires
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y-32.0*10)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*10)
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y-32.0*11)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*11)
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y-32.0*12)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*12)
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y-32.0*13)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*13)
	w.Image.DrawFrame(7, w.Rect.X, w.Rect.Y-32.0*14)
	w.Image.DrawFrame(8, w.Rect.X+32.0, w.Rect.Y-32.0*14)
}

// Bounds TODO doc
func (w *Wall) Bounds() shapes.Rect {
	return *(w.Rect)
}
