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

package app

import (
	"fmt"
	_ "image/png" // Need this for image libs
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/FlappyDisk/actors/mountains"
	"github.com/hurricanerix/FlappyDisk/actors/player"
	"github.com/hurricanerix/FlappyDisk/input"
	"github.com/hurricanerix/FlappyDisk/window"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Config contains settings for running the app
type Config struct {
	Window window.Config
	Input  input.Config
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
func (a Config) Run() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var monitor *glfw.Monitor
	if a.Window.FullScreen {
		// TODO: Maybe choose monitor based on config?
		// http://www.glfw.org/docs/latest/monitor.html#monitor_monitors
		monitor = glfw.GetPrimaryMonitor()
	}

	window, err := glfw.CreateWindow(a.Window.Width, a.Window.Height, "Flappy Disk", monitor, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	window.SetKeyCallback(keyCallback)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	var left, right, top, bottom, near, far float32
	right = float32(a.Window.Width)
	top = float32(a.Window.Height)
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

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		player.Update(elapsed)
		mountains.Update(elapsed)

		if player.Dead {
			fmt.Println("You died!")
			window.SetShouldClose(true)
		}

		// Render
		gl.UseProgram(program)

		// Drawing order matters here, draw from back to front.
		mountains.Draw()
		player.Draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShader = `
#version 330

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;
uniform mat4 ModelMatrix;

in vec3 MCVertex;
in vec2 TexCoord0;

out vec2 TexCoord;
out float Layer;

void main() {
  TexCoord = TexCoord0;
  gl_Position = ProjMatrix * ViewMatrix * ModelMatrix * vec4(MCVertex, 1);
	//gl_Position.z = MCVertex.z;
	Layer = gl_Position.z;
}
` + "\x00"

var fragmentShader = `
#version 330

uniform sampler2D ColorMap;

in vec2 TexCoord;
in float Layer;

out vec4 outputColor;

void main() {
	outputColor = texture(ColorMap, TexCoord);
}
` + "\x00"
