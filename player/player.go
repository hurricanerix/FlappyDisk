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
// Package player TODO doc

package player

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/shapes"
	"github.com/hurricanerix/shade/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Player struct {
	Sprite  *sprite.Context
	Rect    *shapes.Rect
	Alive   bool
	dy      float32
	jumpKey bool
}

// New TODO doc
func New(x, y float32, s *sprite.Context, group *sprite.Group) (*Player, error) {
	// TODO should take a group in as a argument
	p := Player{
		Sprite: s,
		Alive:  true,
	}

	rect, err := shapes.NewRect(x, y, float32(p.Sprite.Width), float32(p.Sprite.Height))
	if err != nil {
		return &p, fmt.Errorf("could create rect: %v", err)
	}
	p.Rect = rect

	// TODO: this should probably be added outside of player
	group.Add(&p)
	return &p, nil
}

// HandleEvent TODO doc
func (p *Player) HandleEvent(event events.Event, dt float32) {
	// TODO: move this to SDK to handle things like holding Left & Right at the same time correctly

	if (event.Action == glfw.Press || event.Action == glfw.Repeat) && event.Key == glfw.KeySpace {
		p.jumpKey = true
	}
	if event.Action == glfw.Release && event.Key == glfw.KeySpace {
		p.jumpKey = false
	}
}

// Bind TODO doc
func (p *Player) Bind(program uint32) error {
	return p.Sprite.Bind(program)
}

// Update TODO doc
func (p *Player) Update(dt float32, g *sprite.Group) {
	if p.jumpKey {
		p.dy = 500.0
	}
	p.dy = float32(math.Max(float64(-400.0), float64(p.dy-40.0)))

	p.Rect.Y += p.dy * dt
	if p.Rect.Top() < 0 {
		p.Alive = false
		p.Rect.Y = 0.0 - float32(p.Sprite.Height)
	}
	if p.Rect.Bottom() > 480 {
		p.Rect.Y = 480
	}

	for _, cell := range sprite.Collide(p, g, false) {
		if cell != nil {
			p.Alive = false
		}
	}
}

// Draw TODO doc
func (p *Player) Draw() {
	p.Sprite.Draw(mgl32.Vec3{p.Rect.X, p.Rect.Y, 0}, nil)
}

// Bounds TODO doc
func (p *Player) Bounds() chan shapes.Rect {
	b := make(chan shapes.Rect, 1)
	b <- *p.Rect
	close(b)
	return b
}
