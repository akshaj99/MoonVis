package render

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Sphere struct {
	vao        uint32
	vbo        uint32
	ebo        uint32
	indices    []uint32
	vertices   []float32
	numIndices int32
}

func NewSphere(radius float32, sectors, stacks int) *Sphere {
	var vertices []float32
	var indices []uint32

	// Create vertices
	for i := 0; i <= stacks; i++ {
		v := float32(i) / float32(stacks)
		phi := float64(math.Pi * float64(i) / float64(stacks))

		for j := 0; j <= sectors; j++ {
			u := float32(j) / float32(sectors)
			theta := float64(2 * math.Pi * float64(j) / float64(sectors))

			// Calculate vertex position
			x := float32(math.Cos(theta) * math.Sin(phi))
			y := float32(math.Cos(phi))
			z := float32(math.Sin(theta) * math.Sin(phi))

			// Add position (scaled by radius)
			vertices = append(vertices, x*radius)
			vertices = append(vertices, y*radius)
			vertices = append(vertices, z*radius)

			// Add normal (same as position for unit sphere)
			vertices = append(vertices, x)
			vertices = append(vertices, y)
			vertices = append(vertices, z)

			// Add texture coordinates
			vertices = append(vertices, u)
			vertices = append(vertices, v)
		}
	}

	// Generate indices
	for i := 0; i < stacks; i++ {
		k1 := i * (sectors + 1)
		k2 := k1 + sectors + 1

		for j := 0; j < sectors; j++ {
			if i != 0 {
				indices = append(indices, uint32(k1), uint32(k2), uint32(k1+1))
			}
			if i != (stacks - 1) {
				indices = append(indices, uint32(k1+1), uint32(k2), uint32(k2+1))
			}
			k1++
			k2++
		}
	}

	// Create OpenGL buffers
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Normal attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Texture coordinate attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	return &Sphere{
		vao:        vao,
		vbo:        vbo,
		ebo:        ebo,
		indices:    indices,
		vertices:   vertices,
		numIndices: int32(len(indices)),
	}
}

func (s *Sphere) Draw() {
	gl.BindVertexArray(s.vao)
	gl.DrawElements(gl.TRIANGLES, s.numIndices, gl.UNSIGNED_INT, nil)
}

func (s *Sphere) Cleanup() {
	gl.DeleteVertexArrays(1, &s.vao)
	gl.DeleteBuffers(1, &s.vbo)
	gl.DeleteBuffers(1, &s.ebo)
}
