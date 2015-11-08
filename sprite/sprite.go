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
	}

	return &s, nil
}

// Sprite represents position, rotation and scale for a given asset.
type Sprite struct {
	AssetName string
	Pos       mgl32.Vec3
	Scale     float32
	Rot       float32
	Texture   uint32
	data      []byte
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
