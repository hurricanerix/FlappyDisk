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

func getView(scaleX, scaleY, transX, transY, transZ float32) mgl32.Mat4 {
	v := mgl32.Mat4{
		scaleX, 0.0, 0.0, transX,
		0.0, scaleY, 0.0, transY,
		0.0, 0.0, 1.0, transZ,
		0.0, 0.0, 0.0, 1.0,
	}
	return v
}

// getProj
//
// U{Modern glOrtho2d<http://stackoverflow.com/questions/21323743/
//   modern-equivalent-of-gluortho2d>}
//
//  U{Orthographic Projection<http://en.wikipedia.org/wiki/
//    Orthographic_projection_(geometry)>}
//
//  @param left: position of the left side of the display
//  @type left: int
//  @param right: position of the right side of the display
//  @type right: int
//  @param bottom: position of the bottom side of the display
//  @type bottom: int
//  @param top: position of the top side of the display
//  @type top: int
func getProj(left, right, bottom, top float32) mgl32.Mat4 {
	/*
	   mat = [
	       (2.0 * inv_x), 0.0, 0.0, (-(right + left) * inv_x),
	       0.0, (2.0 * inv_y), 0.0, (-(top + bottom) * inv_y),
	       0.0, 0.0, (-2.0 * inv_z), (-(zFar + zNear) * inv_z),
	       0.0, 0.0, 0.0, 1.0]
	*/
	zNear := -25.0
	zFar := 25.0
	invZ := 1.0 / (zFar - zNear)
	invY := 1.0 / (top - bottom)
	invX := 1.0 / (right - left)

	m := mgl32.Mat4{
		(2.0 * invX), 0.0, 0.0, (-(right + left) * invX),
		0.0, (2.0 * invY), 0.0, (-(top + bottom) * invY),
		0.0, 0.0, float32(-2.0 * invZ), float32(-(zFar + zNear) * invZ),
		0.0, 0.0, 0.0, 1.0,
	}
	return m
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

	var monitor *glfw.Monitor = nil
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

	// Configure global settings
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.527, 0.805, 0.918, 1.0)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	//left := 0.0
	right := a.Window.Width
	//top := 0.0
	bottom := a.Window.Height
	//near := -25.0
	//far := 25.0

	//projection := mgl32.Ortho(float32(left), float32(right), float32(bottom), float32(top), float32(near), float32(far))
	projection := getProj(1.0, float32(right), 0.0, float32(bottom))
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	//camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0})
	camera := getView(1.0, float32(right), 0.0, float32(bottom), 1.0)
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

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

	//angle := 0.0
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

		player.Draw()
		mountains.Draw()

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

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
