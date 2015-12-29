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
// Package sprite provides functions for managing sprites.
// NOTE:

package sprite

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/png" // register PNG decode
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hurricanerix/FlappyDisk/gen"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// New returns Sprite
func New(assetName string) (*Sprite, error) {
	data, err := gen.Asset(assetName)
	if err != nil {
		return nil, err
	}

	tex, width, height, err := newTexture(data)
	if err != nil {
		return nil, err
	}

	s := Sprite{
		AssetName:    assetName,
		CurrentFrame: 0,
		FrameCount:   int(width / 32.0),
		Width:        float32(width),
		Height:       float32(height),
		Texture:      tex,
		data:         data,
		model:        mgl32.Ident4(),
	}

	return &s, nil
}

// Name
// SheetWidth
// FrameHeight
// FrameWidth
// FrameCount
// data

// Sprite represents position, rotation and scale for a given asset.
type Sprite struct {
	AssetName    string // Name
	CurrentFrame int
	FrameCount   int
	Width        float32
	Height       float32
	Texture      uint32
	data         []byte
	vao          uint32
	frame        mgl32.Mat3
	frameMatrix  int32
	model        mgl32.Mat4
	modelMatrix  int32
}

// Bind TODO: write comment
func (s *Sprite) Bind(program uint32) error {
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	mcVertex := uint32(gl.GetAttribLocation(program, gl.Str("MCVertex\x00")))
	gl.EnableVertexAttribArray(mcVertex)
	gl.VertexAttribPointer(mcVertex, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoord0 := uint32(gl.GetAttribLocation(program, gl.Str("TexCoord0\x00")))
	gl.EnableVertexAttribArray(texCoord0)
	gl.VertexAttribPointer(texCoord0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	colorMap := gl.GetUniformLocation(program, gl.Str("ColorMap\x00"))
	gl.Uniform1i(colorMap, 0)

	s.frameMatrix = gl.GetUniformLocation(program, gl.Str("FrameMatrix\x00"))
	gl.UniformMatrix4fv(s.frameMatrix, 1, false, &s.frame[0])

	s.modelMatrix = gl.GetUniformLocation(program, gl.Str("ModelMatrix\x00"))
	gl.UniformMatrix4fv(s.modelMatrix, 1, false, &s.model[0])
	return nil
}

// Draw TODO: write comment
func (s *Sprite) Draw(rotation float32, translation mgl32.Vec3, scale float32) {
	//gl.Enable(gl.DEPTH_TEST)

	// Calculate Frame
	texScale := 1.0 / float32(s.FrameCount)
	//texTrans := float32(s.CurrentFrame) * texScale
	texTrans := texScale //float32(s.CurrentFrame) * texScale

	// scale_x = 1.0/self.data['frame']['count']['x']
	// scale_y = 1.0/self.data['frame']['count']['y']
	// trans_x = frame_x * scale_x
	// trans_y = frame_y * scale_y

	// def get_3x3_transform(scale_x=1.0, scale_y=1.0, trans_x=1.0, trans_y=1.0):
	//     """Returns a 3x3 transform.
	//
	//     @return: transformation matrix
	//     @rtype: list
	//     """
	//     transform = [scale_x, 0.0, trans_x,
	//                  0.0, scale_y, trans_y,
	//                  0.0, 0.0, 1.0]
	//     return transform

	s.frame = mgl32.Mat3{
		float32(texScale), 0.0, float32(texTrans),
		0.0, 1.0, 1.0,
		0.0, 0.0, 1.0,
	}

	// Calculate Model
	s.model = mgl32.Ident4()
	s.model = s.model.Mul4(mgl32.Translate3D(translation.X(), translation.Y(), translation.Z()))
	s.model = s.model.Mul4(mgl32.HomogRotate3DZ(rotation))
	s.model = s.model.Mul4(mgl32.Scale3D(32.0*scale, 32.0*scale, 1.0))

	gl.UniformMatrix3fv(s.frameMatrix, 1, false, &s.frame[0])
	gl.UniformMatrix4fv(s.modelMatrix, 1, false, &s.model[0])
	gl.BindVertexArray(s.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.Texture)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

func newTexture(b []byte) (uint32, int, int, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return 0, 0, 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, 0, 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32

	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, rgba.Rect.Size().X, rgba.Rect.Size().Y, nil
}

var vertices = []float32{
	-1.0, -1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
}

var VertexShader = `
#version 330

uniform mat4 ProjMatrix;
uniform mat4 ViewMatrix;
uniform mat4 ModelMatrix;
uniform mat3 FrameMatrix;

in vec3 MCVertex;
in vec2 TexCoord0;

out vec2 TexCoord;
out float Layer;

void main() {
	TexCoord = vec3(FrameMatrix * vec3(TexCoord0, 0.0)).st;
  gl_Position = ProjMatrix * ViewMatrix * ModelMatrix * vec4(MCVertex, 1);
	Layer = gl_Position.z;
}
` + "\x00"

var FragmentShader = `
#version 330

uniform sampler2D ColorMap;

in vec2 TexCoord;
in float Layer;

out vec4 outputColor;

void main() {
	outputColor = texture(ColorMap, TexCoord);
}
` + "\x00"
