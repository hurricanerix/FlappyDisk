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

	"github.com/hurricanerix/FlappyDisk/sprite"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func New() (*Mountains, error) {
	s := []*sprite.Sprite{}

	for i := 0; i < 2; i++ {
		si, err := sprite.New("assets/mountains.png")
		if err != nil {
			return nil, err
		}
		s = append(s, si)
	}

	m := Mountains{
		Sprite: s,
	}

	return &m, nil

}

type Mountains struct {
	Sprite []*sprite.Sprite
}

func (m *Mountains) Bind(program uint32) {
	for i := 0; i < len(m.Sprite); i++ {
		m.Sprite[i].Bind(program)
	}
}

func (m *Mountains) Update(elapsed float64) {
	m.Sprite[0].Rot -= elapsed
}

func (m *Mountains) Draw() {
	for i := 0; i < len(m.Sprite); i++ {
		m.Sprite[i].Draw()
	}
}