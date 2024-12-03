package main

import (
	"fmt"
	"math"
	"moonVis/pkg/loader"
	"moonVis/pkg/render"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 1280
	height = 720
)

var (
	cameraPos              = mgl32.Vec3{0, 0, 3.5}
	cameraFront            = mgl32.Vec3{0, 0, -1}
	cameraUp               = mgl32.Vec3{0, 1, 0}
	rotationMatrix         = mgl32.Ident4()
	isDragging             = false
	lastX          float64 = 0.0
	lastY          float64 = 0.0
)

func main() {
	window, err := render.NewWindow()
	if err != nil {
		fmt.Printf("Error creating window: %v\n", err)
		return
	}
	defer window.Terminate()

	// Load textures
	colorTexture, err := loader.LoadTexture("data/lroc_color_poles_8k.tif")
	if err != nil {
		fmt.Printf("Error loading color texture: %v\n", err)
		return
	}

	shader, err := render.NewShader()
	if err != nil {
		fmt.Printf("Error creating shader: %v\n", err)
		return
	}

	// Create sphere
	sphere := render.NewSphere(1.0, 360, 180)
	defer sphere.Cleanup()

	// Set up mouse callback
	window.SetMouseCallback(mouseCallback)
	window.SetScrollCallback(scrollCallback)

	gl.Enable(gl.DEPTH_TEST)

	// Camera and projection setup
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 100.0)

	lastTime := glfw.GetTime()

	for !window.ShouldClose() {
		currentTime := glfw.GetTime()
		deltaTime := float32(currentTime - lastTime)
		lastTime = currentTime

		processInput(window.GlfwWindow(), deltaTime)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)

		model := rotationMatrix

		view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)

		shader.Use()

		modelLoc := gl.GetUniformLocation(shader.Program, gl.Str("model\x00"))
		viewLoc := gl.GetUniformLocation(shader.Program, gl.Str("view\x00"))
		projLoc := gl.GetUniformLocation(shader.Program, gl.Str("projection\x00"))

		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
		gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, colorTexture)
		gl.Uniform1i(gl.GetUniformLocation(shader.Program, gl.Str("moonTexture\x00")), 0)

		sphere.Draw()

		window.SwapBuffers()
		window.PollEvents()
	}
}

func sphericalToCartesian(lat, lon float32) mgl32.Vec3 {
	latRad := mgl32.DegToRad(lat)
	lonRad := mgl32.DegToRad(lon)

	cosLat := float32(math.Cos(float64(latRad)))
	sinLat := float32(math.Sin(float64(latRad)))
	cosLon := float32(math.Cos(float64(lonRad)))
	sinLon := float32(math.Sin(float64(lonRad)))

	return mgl32.Vec3{
		cosLat * cosLon,
		sinLat,
		cosLat * sinLon,
	}
}

func projectToScreen(pos mgl32.Vec3, view, projection mgl32.Mat4) mgl32.Vec2 {
	// Transform to clip space
	clipPos := projection.Mul4(view).Mul4x1(pos.Vec4(1.0))

	if clipPos.W() != 0 {
		clipPos = clipPos.Mul(1 / clipPos.W())
	}

	// Transform to screen coordinates
	screenX := (clipPos.X() + 1) * float32(width) / 2
	screenY := (1 - clipPos.Y()) * float32(height) / 2

	return mgl32.Vec2{screenX, screenY}
}

func mouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	if w.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		if !isDragging {
			isDragging = true
			lastX = xpos
			lastY = ypos
			return
		}

		// Calculate rotation based on mouse movement
		xoffset := float32(xpos - lastX)
		yoffset := float32(ypos - lastY)
		lastX = xpos
		lastY = ypos

		sensitivity := float32(0.005)

		// Create rotation matrices for X and Y axes
		rotateY := mgl32.HomogRotate3D(xoffset*sensitivity, mgl32.Vec3{0, 1, 0})
		rotateX := mgl32.HomogRotate3D(yoffset*sensitivity, mgl32.Vec3{1, 0, 0})

		// Apply rotations to global rotation matrix
		rotationMatrix = rotateY.Mul4(rotateX).Mul4(rotationMatrix)
	} else {
		isDragging = false
	}
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	zoomSpeed := float32(0.1)
	newPos := cameraPos.Add(cameraFront.Mul(float32(yoff) * zoomSpeed))
	distance := newPos.Len()
	if distance > 1.4 && distance < 10.0 {
		cameraPos = newPos
	}
}

func processInput(window *glfw.Window, deltaTime float32) {
	speed := float32(2.5) * deltaTime
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Mul(speed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Mul(speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(speed))
	}
}
