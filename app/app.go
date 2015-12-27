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
// Package app manages the main game loop.

package app

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"gopkg.in/gcfg.v1"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/FlappyDisk/actors/mountains"
	"github.com/hurricanerix/FlappyDisk/actors/player"
	"github.com/hurricanerix/FlappyDisk/gen"
	"github.com/hurricanerix/FlappyDisk/input"
	"github.com/hurricanerix/FlappyDisk/shader"
	"github.com/hurricanerix/FlappyDisk/sprite"
	"github.com/hurricanerix/FlappyDisk/window"
)

// Config of the appliction
type Config struct {
	Window window.Config
	Input  input.Config
}

// Context of the application.
type Context struct {
	// Config of application
	Config Config

	// Monitor (GLFW) to use for fullscreen, or nil for windowed.
	Monitor *glfw.Monitor

	// Window (GLFW) context.
	Window *glfw.Window

	// Program (GLSL) used for rendering the scene.
	Program uint32
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// New Context for the application.
func New(resetConf bool) (*Context, error) {

	configPath, configName := getConfigPathName()

	if resetConf {
		fmt.Println("resetting config to defaults")
		err := createConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var c Config
	err := gcfg.ReadFileInto(&c, configPath+configName)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			createConfig()
		} else {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	// TODO: Verify config settings are valid.

	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var monitor *glfw.Monitor
	if c.Window.FullScreen {
		// TODO: Maybe choose monitor based on config?
		// http://www.glfw.org/docs/latest/monitor.html#monitor_monitors
		monitor = glfw.GetPrimaryMonitor()
	}

	window, err := glfw.CreateWindow(c.Window.Width, c.Window.Height, "Flappy Disk", monitor, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		return nil, err
	}

	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	window.SetKeyCallback(keyCallback)

	context := Context{
		Config:  c,
		Monitor: monitor,
		Window:  window,
	}

	return &context, nil
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// TODO: Read from config
	if key == glfw.KeyBackspace && action == glfw.Press {
		println("Select")
	}
	if key == glfw.KeyEnter && action == glfw.Press {
		println("Start")
	}
	if key == glfw.KeySpace && action == glfw.Press {
		println("Flap")
	}
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

// Run the application
func (app Context) Run() {

	// Configure the vertex and fragment shaders
	program, err := shader.New(sprite.VertexShader, sprite.FragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	var left, right, top, bottom, near, far float32
	right = float32(app.Config.Window.Width)
	top = float32(app.Config.Window.Height)
	near = 0.1
	far = 10.0

	projMatrix := mgl32.Ortho(left, right, bottom, top, near, far)
	projUniform := gl.GetUniformLocation(program, gl.Str("ProjMatrix\x00"))
	gl.UniformMatrix4fv(projUniform, 1, false, &projMatrix[0])

	var eye, center, up mgl32.Vec3
	eye = mgl32.Vec3{0.0, 0.0, 7.0}
	center = mgl32.Vec3{0.0, 0.0, -1.0}
	up = mgl32.Vec3{0.0, 1.0, 0.0}
	viewMatrix := mgl32.LookAtV(eye, center, up)

	viewUniform := gl.GetUniformLocation(program, gl.Str("ViewMatrix\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &viewMatrix[0])

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	player, err := player.New()
	if err != nil {
		panic(err)
	}

	mountains, err := mountains.New()
	if err != nil {
		panic(err)
	}

	player.Bind(program)
	mountains.Bind(program)

	// Configure global settings
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	// TODO: Figure out why "layering" using z-buffer does not work.
	//gl.DepthFunc(gl.LESS)
	//gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.NEVER)
	gl.Enable(gl.CULL_FACE)

	gl.ClearColor(0.527, 0.805, 0.918, 1.0)

	previousTime := glfw.GetTime()

	for !app.Window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		player.Update(elapsed)
		mountains.Update(elapsed)

		if player.Dead {
			fmt.Println("You died!")
			app.Window.SetShouldClose(true)
		}

		// Render
		gl.UseProgram(program)

		// Drawing order matters here, draw from back to front.
		mountains.Draw()
		player.Draw()

		// Maintenance
		app.Window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (app Context) Terminate() {
	glfw.Terminate()
}

func getConfigPathName() (string, string) {
	usr, _ := user.Current()
	return usr.HomeDir + "/.config/flappy-disk/", "app.conf"
}

func createConfig() error {
	path, name := getConfigPathName()
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(path + name)
	defer f.Close()
	if err != nil {
		return err
	}

	configData, err := gen.Asset("assets/default.conf")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = f.Write(configData)
	if err != nil {
		return err
	}

	f.Sync()

	return nil
}
