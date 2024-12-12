package render

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	Program uint32
}

//go:embed shaders/vertex.glsl
var vertexShaderSource string

//go:embed shaders/fragment.glsl
var fragmentShaderSource string

func NewShader() (*Shader, error) {
	//Add null terminator to shader sources
	vertexSource := vertexShaderSource + "\x00"
	fragmentSource := fragmentShaderSource + "\x00"

	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return nil, fmt.Errorf("failed to compile vertex shader: %v", err)
	}

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, fmt.Errorf("failed to compile fragment shader: %v", err)
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
		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{Program: program}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.Program)
}

// Added helper methods for setting uniforms
func (s *Shader) SetBool(name string, value bool) {
	var intValue int32
	if value {
		intValue = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), intValue)
}

func (s *Shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetVec3(name string, value mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) SetMat4(name string, value mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), 1, false, &value[0])
}
