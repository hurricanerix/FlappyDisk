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
// Package game manages the main game loop.

package game

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/FlappyDisk/player"
	"github.com/hurricanerix/FlappyDisk/walls"
	"github.com/hurricanerix/shade/display"
	"github.com/hurricanerix/shade/events"
	"github.com/hurricanerix/shade/fonts"
	"github.com/hurricanerix/shade/sprite"
	"github.com/hurricanerix/shade/time/clock"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Config TODO doc
type Config struct {
	Cheat bool
}

// Context TODO doc
type Context struct {
	Screen *display.Context
	Player *player.Player
	Walls  *sprite.Group
}

// New TODO doc
func New(screen *display.Context) (Context, error) {
	return Context{
		Screen: screen,
	}, nil
}

// Main TODO doc
func (c *Context) Main(screen *display.Context, config Config) {
	clock, err := clock.New()
	if err != nil {
		panic(err)
	}

	//background, err := sprite.Load("background.png", 1)
	//if err != nil {
	//	panic(err)
	//}
	//background.Bind(c.Screen.Program)

	sprites := sprite.NewGroup()
	p, err := player.New(sprites)
	if err != nil {
		panic(err)
	}

	c.Walls = sprite.NewGroup()
	sprites.Add(c.Walls)

	_, err = walls.New(c.Walls)
	if err != nil {
		panic(err)
	}
	//_, err = walls.New(false, 120, c.Walls)
	//if err != nil {
	//panic(err)
	//}

	// TODO: should only load image data once.
	//block, err := sprite.Load("transistor.png", 1)
	//if err != nil {
	//	panic(err)
	//}
	//println(block)
	font, err := fonts.New()
	if err != nil {
		panic(err)
	}
	font.Bind(screen.Program)

	sprites.Bind(c.Screen.Program)
	for running := true; running; {
		p.Alive = true
		dt := clock.Tick(30)

		// TODO move this somewhere else (maybe a Clear method of display
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// TODO refector events to be cleaner
		if screen.Window.ShouldClose() {
			running = !screen.Window.ShouldClose()
		}

		for _, event := range events.Get() {
			if event.Action == glfw.Press && event.Key == glfw.KeyEscape {
				running = false
				event.Window.SetShouldClose(true)
			}
			p.HandleEvent(event, dt/1000.0)
		}

		sprites.Update(dt/1000.0, c.Walls)
		screen.Fill(200.0/256.0, 200/256.0, 200/256.0)
		//background.Draw(0, 0)
		if p.Alive == false {
			msg := "You Died!"
			font.DrawText(250, 250, 2.0, 2.0, nil, msg)
			if !config.Cheat {
				running = false
			}
		}

		sprites.Draw()

		// TODO: implement score
		msg := fmt.Sprintf("%d", 0)
		w, h := font.SizeText(3.0, 3.0, msg)
		font.DrawText(screen.Width/2-w/2, screen.Height-h, 3.0, 3.0, nil, msg)

		if config.Cheat {
			msg := "Dev Mode!\n"
			msg += fmt.Sprintf("Pos: %.0f, %.0f\n", p.Rect.X, p.Rect.Y)
			msg += fmt.Sprintf("Status: %t\n", p.Alive)
			_, h := font.SizeText(1.0, 1.0, msg)
			font.DrawText(0, 480-h, 2.0, 2.0, &mgl32.Vec4{0.0, 0.0, 0.0, 0.5}, msg)
		}
		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}
}
