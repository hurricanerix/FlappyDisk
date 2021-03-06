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
	"bytes"
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/FlappyDisk/gen"
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
	c.Walls = sprite.NewGroup()

	playerSprite, err := loadSprite("assets/floppy.png", 1, 1)
	if err != nil {
		panic(err)
	}
	p, err := player.New(320.0, 240.0, playerSprite, sprites)
	if err != nil {
		panic(err)
	}

	wallSprite, err := loadSprite("assets/resistor.png", 32, 1)
	if err != nil {
		panic(err)
	}

	sprites.Add(c.Walls)

	_, err = walls.New(640, 240, 80, wallSprite, c.Walls)
	if err != nil {
		panic(err)
	}

	font, err := fonts.SimpleASCII()
	if err != nil {
		panic(err)
	}
	font.Bind(screen.Program)

	sprites.Bind(c.Screen.Program)
	for running := true; running; {
		screen.Fill(200.0/256.0, 200/256.0, 200/256.0)

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

		//background.Draw(0, 0)
		if p.Alive == false {
			msg := "You Died!"
			font.DrawText(mgl32.Vec3{250, 250, 0}, &sprite.Effects{Scale: mgl32.Vec3{2.0, 2.0, 1.0}}, msg)
			if !config.Cheat {
				running = false
			}
		}

		sprites.Draw(nil)

		// TODO: implement score
		msg := fmt.Sprintf("%d", 0)
		effect := sprite.Effects{Scale: mgl32.Vec3{3.0, 3.0, 1.0}}
		w, h := font.SizeText(&effect, msg)
		font.DrawText(mgl32.Vec3{screen.Width/2 - w/2, screen.Height - h, 0}, &effect, msg)

		if config.Cheat {
			msg := "Dev Mode!\n"
			msg += fmt.Sprintf("Pos: %.0f, %.0f\n", p.Rect.X, p.Rect.Y)
			msg += fmt.Sprintf("Status: %t\n", p.Alive)
			_, h := font.SizeText(nil, msg)
			font.DrawText(mgl32.Vec3{0, 480 - h, 0}, nil, msg)
		}
		screen.Flip()

		// TODO refector events to be cleaner
		glfw.PollEvents()
	}
}

func loadSprite(name string, framesWide, framesHigh int) (*sprite.Context, error) {
	imgFile, err := gen.Asset(name)
	if err != nil {
		return nil, fmt.Errorf("could not load asset %s: %v", name, err)
	}
	i, _, err := image.Decode(bytes.NewReader(imgFile))
	if err != nil {
		return nil, fmt.Errorf("could not decode file %s: %v", name, err)
	}

	s, err := sprite.New(i, nil, framesWide, framesHigh)
	if err != nil {
		return nil, err
	}

	return s, nil
}
