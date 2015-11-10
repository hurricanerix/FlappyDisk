// Copyright 2015 Richard Hawkins
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
// Package app provides the starting point for the app.

package player

import (
	"runtime"

	"github.com/hurricanerix/FlappyDisk/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func New() (*Player, error) {
	s, err := sprite.New("assets/floppy.png")
	if err != nil {
		return nil, err
	}

	p := Player{
		Sprite:  s,
		Falling: true,
		Dead:    false,
	}

	return &p, nil

}

type Player struct {
	Sprite  *sprite.Sprite
	Falling bool
	Dead    bool
}

func (p *Player) Bind(program uint32) {
	p.Sprite.Bind(program)
}

func (p *Player) Update(elapsed float64) {
	p.Sprite.Rot -= (elapsed * 2)
	p.Sprite.Scale = 0.5
	if p.Falling {
		p.Sprite.Pos[1] -= float32(elapsed) * 3
	} else {
		p.Sprite.Pos[1] += float32(elapsed) * 3
	}
	if p.Sprite.Pos[1] < -3 {
		p.Dead = true
	}
	p.Sprite.Pos[2] = 10
}

func (p *Player) Draw() {
	p.Sprite.Draw()
}
