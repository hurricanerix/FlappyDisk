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

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/hurricanerix/FlappyDisk/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func New(window *glfw.Window) (*Player, error) {
	s, err := sprite.New("assets/floppy.png")
	if err != nil {
		return nil, err
	}

	p := Player{
		Sprite: s,
	}

	//window.SetKeyCallback(flapCallback)

	return &p, nil

}

type Player struct {
	Sprite *sprite.Sprite
}

func (p *Player) Bind(program uint32) {
	p.Sprite.Bind(program)
}

func (p *Player) Update(elapsed float64) {
	p.Sprite.Rot += elapsed
}

func (p *Player) Draw() {
	p.Sprite.Draw()
}
