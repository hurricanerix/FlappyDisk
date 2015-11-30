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

package mountains

import (
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/FlappyDisk/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// New TODO: write comment
func New() (*Mountains, error) {
	s, err := sprite.New("assets/mountains.png")
	if err != nil {
		return nil, err
	}

	m := Mountains{
		Sprite: s,
		Pos:    mgl32.Vec3{0.0, 0.0, 2.0},
	}

	return &m, nil
}

// Mountains TODO: write comment
type Mountains struct {
	Sprite *sprite.Sprite
	Pos    mgl32.Vec3
}

// Bind TODO: write comment
func (m *Mountains) Bind(program uint32) {
	m.Sprite.Bind(program)
}

// Update TODO: write comment
func (m *Mountains) Update(elapsed float64) {
}

// Draw TODO: write comment
func (m *Mountains) Draw() {
	m.Sprite.Draw(0.0, m.Pos, 1.0)
}
