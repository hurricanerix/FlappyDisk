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
	"github.com/hurricanerix/transylvania/events"
	"github.com/hurricanerix/transylvania/shapes"
	"github.com/hurricanerix/transylvania/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Player TODO doc
type Player struct {
	Image   *sprite.Context
	Rect    *shapes.Rect
	Alive   bool
	dy      float32
	jumpKey bool
}

// New TODO doc
func New(group *sprite.Group) (*Player, error) {
	// TODO should take a group in as a argument
	p := Player{
		Alive: true,
	}

	player, err := sprite.Load("floppy.png")
	if err != nil {
		return &p, fmt.Errorf("could not load player: %v", err)
	}
	p.Image = player

	rect, err := shapes.NewRect(320.0, 240.0, float32(p.Image.Width), float32(p.Image.Height))
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
	return p.Image.Bind(program)
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
	}

}

// Draw TODO doc
func (p *Player) Draw() {
	p.Image.Draw(p.Rect.X, p.Rect.Y)
}

// Bounds TODO doc
func (p *Player) Bounds() shapes.Rect {
	return *(p.Rect)
}
