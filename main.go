// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// An app that draws a green triangle on a red background.
package main

import (
	"encoding/binary"
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

var (
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer
	boxColor Color
	opacity  float32
	touchLoc geom.Point
	center   geom.Point
)

func main() {
	app.Run(app.Callbacks{
		Draw:  draw,
		Touch: touch,
	})
}

// TODO(crawshaw): Need an easier way to do GL-dependent initialization.
type Color struct {
	Red, Green, Blue float32
}

func setColor(r, g, b int16) Color {
	return Color{
		float32(r) / 255,
		float32(g) / 255,
		float32(b) / 255,
	}
}

func initGL() {
	var err error
	program, err = glutil.CreateProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = gl.GenBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, gl.STATIC_DRAW, scene)

	position = gl.GetAttribLocation(program, "position")
	color = gl.GetUniformLocation(program, "color")
	offset = gl.GetUniformLocation(program, "offset")
	center = geom.Point{geom.Width / 2, geom.Height / 2}
	touchLoc = center
	boxColor = setColor(192, 192, 192)
}

func touch(t event.Touch) {
	touchLoc = t.Loc
}

func draw() {
	if program.Value == 0 {
		initGL()
	}

	gl.ClearColor(0, 0, 0, 0.01)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)

	gl.Uniform4f(color, boxColor.Red, boxColor.Green, boxColor.Blue, 0.01)

	gl.Uniform2f(offset, float32(touchLoc.X/geom.Width), float32(touchLoc.Y/geom.Height))

	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.EnableVertexAttribArray(position)
	gl.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	gl.DrawArrays(gl.TRIANGLES, 1, vertexCount)
	gl.DisableVertexAttribArray(position)
}

var scene = f32.Bytes(binary.LittleEndian,
	0.0, 0.0, 0.0, // bottom left
	0.4, 0.0, 0.0, // bottom right
	0.0, 0.4, 0.0, // top left
	0.4, 0.4, 0.0, // top right
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
