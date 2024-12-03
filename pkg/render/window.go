package render

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 1280
	height = 720
)

type Window struct {
	window *glfw.Window
}

func (w *Window) SetMouseCallback(callback func(w *glfw.Window, xpos float64, ypos float64)) {
	w.window.SetCursorPosCallback(callback)
}

func (w *Window) SetScrollCallback(callback func(w *glfw.Window, xoff float64, yoff float64)) {
	w.window.SetScrollCallback(callback)
}

func (w *Window) GlfwWindow() *glfw.Window {
	return w.window
}

func init() {
	runtime.LockOSThread()
}

func NewWindow() (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %v", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Moon Visualization", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %v", err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize OpenGL: %v", err)
	}

	return &Window{window: window}, nil
}

func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w *Window) SwapBuffers() {
	w.window.SwapBuffers()
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) Terminate() {
	glfw.Terminate()
}
