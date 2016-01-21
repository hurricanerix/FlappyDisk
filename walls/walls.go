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
	"math/rand"
	"runtime"

	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Wall struct {
	Image      *sprite.Context
	TopRect    *shapes.Rect
	BottomRect *shapes.Rect
	width      float32
	dx         float32
	offset     float32
	size       float32
}

// New TODO doc
func New(group *sprite.Group) (*Wall, error) {
	// TODO should take a group in as a argument
	w := Wall{
		width:  32.0 * 2,
		offset: 240.0,
		size:   80.0,
	}

	wall, err := sprite.Load("resistor.png", 32, 1)
	if err != nil {
		return &w, fmt.Errorf("could not load wall: %v", err)
	}
	w.Image = wall

	topRect, err := shapes.NewRect(640.0-w.width, w.offset+w.size/2.0, 64.0, 480.0)
	if err != nil {
		return &w, fmt.Errorf("could not create top rect: %v", err)
	}
	w.TopRect = topRect

	bottomRect, err := shapes.NewRect(640.0-w.width, w.offset-w.size/2.0, 64.0, 480.0)
	if err != nil {
		return &w, fmt.Errorf("could create bottom rect: %v", err)
	}
	w.BottomRect = bottomRect

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
	w.TopRect.X += w.dx * dt
	w.BottomRect.X += w.dx * dt
	if w.TopRect.X+float32(w.Image.Width) < 0.0 {
		var min int = int(w.size)
		var max int = 480 - int(w.size)
		w.offset = float32(rand.Intn(max-min) + min)
		w.TopRect.X = 641
		w.BottomRect.X = 641
		w.TopRect.Y = w.offset + w.size/2.0
		w.BottomRect.Y = w.offset - w.size/2.0
	}
}

// Draw TODO doc
func (w *Wall) Draw() {
	for i := 32.0 * 3; i < 480; i += 32 {
		w.Image.DrawFrame(7, 0, 1.0, 1.0, w.TopRect.X, w.TopRect.Y+float32(i))
		w.Image.DrawFrame(8, 0, 1.0, 1.0, w.TopRect.X+32.0, w.TopRect.Y+float32(i))
	}
	w.Image.DrawFrame(0, 0, 1.0, 1.0, w.TopRect.X, w.TopRect.Y+32.0*3)
	w.Image.DrawFrame(8, 0, 1.0, 1.0, w.TopRect.X+32.0, w.TopRect.Y+32.0*3)
	w.Image.DrawFrame(1, 0, 1.0, 1.0, w.TopRect.X, w.TopRect.Y+32.0*2)
	w.Image.DrawFrame(8, 0, 1.0, 1.0, w.TopRect.X+32.0, w.TopRect.Y+32.0*2)
	w.Image.DrawFrame(2, 0, 1.0, 1.0, w.TopRect.X, w.TopRect.Y+32.0*1)
	w.Image.DrawFrame(8, 0, 1.0, 1.0, w.TopRect.X+32.0, w.TopRect.Y+32.0*1)
	w.Image.DrawFrame(3, 0, 1.0, 1.0, w.TopRect.X, w.TopRect.Y+32.0*0)
	w.Image.DrawFrame(4, 0, 1.0, 1.0, w.TopRect.X+32, w.TopRect.Y+32.0*0)

	// bottom resistor
	w.Image.DrawFrame(5, 0, 1.0, 1.0, w.BottomRect.X, w.BottomRect.Y-32*1)
	w.Image.DrawFrame(6, 0, 1.0, 1.0, w.BottomRect.X+32.0, w.BottomRect.Y-32.0*1)
	w.Image.DrawFrame(0, 0, 1.0, 1.0, w.BottomRect.X, w.BottomRect.Y-32*2)
	w.Image.DrawFrame(8, 0, 1.0, 1.0, w.BottomRect.X+32.0, w.BottomRect.Y-32*2)
	w.Image.DrawFrame(1, 0, 1.0, 1.0, w.BottomRect.X, w.BottomRect.Y-32.0*3)
	w.Image.DrawFrame(8, 0, 1.0, 1.0, w.BottomRect.X+32.0, w.BottomRect.Y-32.0*3)
	w.Image.DrawFrame(2, 0, 1.0, 1.0, w.BottomRect.X, w.BottomRect.Y-32.0*4)
	w.Image.DrawFrame(8, 0, 1.0, 1.0, w.BottomRect.X+32.0, w.BottomRect.Y-32.0*4)
	for i := 32 * 5; i < 32*13; i += 32 {
		w.Image.DrawFrame(7, 0, 1.0, 1.0, w.BottomRect.X, w.BottomRect.Y-float32(i))
		w.Image.DrawFrame(8, 0, 1.0, 1.0, w.BottomRect.X+32.0, w.BottomRect.Y-float32(i))
	}
}

// Bounds TODO doc
func (w *Wall) Bounds() chan shapes.Rect {
	b := make(chan shapes.Rect, 2)
	b <- *w.TopRect
	b <- *w.BottomRect
	close(b)
	return b
}
