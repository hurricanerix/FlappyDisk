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

package sprite

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
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

	tex, err := newTexture(data)
	if err != nil {
		return nil, err
	}

	s := Sprite{
		AssetName: assetName,
		Pos:       mgl32.Vec3{0.0, 0.0, 0.0},
		Scale:     1.0,
		Rot:       0.0,
		Texture:   tex,
		data:      data,
		model:     mgl32.Ident4(),
	}

	return &s, nil
}

// Sprite represents position, rotation and scale for a given asset.
type Sprite struct {
	AssetName    string
	Pos          mgl32.Vec3
	Scale        float64
	Rot          float64
	Texture      uint32
	data         []byte
	vao          uint32
	model        mgl32.Mat4
	modelUniform int32
}

func (s *Sprite) Bind(program uint32) error {
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	model := mgl32.Ident4()
	s.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(s.modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	return nil
}

func (s *Sprite) Update(elapsed float64) error {
	s.Rot += elapsed
	//model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{1, 1, 1})
	gl.UniformMatrix4fv(s.modelUniform, 1, false, &s.model[0])
	return nil
}

func (s *Sprite) Draw() {

	gl.BindVertexArray(s.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.Texture)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

func newTexture(b []byte) (uint32, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
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

	return texture, nil
}

var vertices = []float32{
	0.0, 0.0, 0.5, 1.0, 0.0,
	0.5, 0.0, 0.5, 0.0, 0.0,
	0.0, 0.5, 0.5, 1.0, 1.0,
	0.5, 0.0, 0.5, 0.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 1.0,
	0.0, 0.5, 0.5, 1.0, 1.0,
}
