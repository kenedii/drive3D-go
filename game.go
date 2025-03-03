package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	Menu GameState = iota
	Playing
)

var (
	currentState    GameState
	prevCarPosition rl.Vector3
)

func initGame() {
	currentState = Menu
	initCar()
	initWorld()
}

func updateGame() {
	switch currentState {
	case Menu:
		handleMenuInput()
	case Playing:
		prevCarPosition = car.position
		updateCar()
		updatePhysics()
		checkCollisions()
	}
}

func drawGame() {
	switch currentState {
	case Menu:
		rl.ClearBackground(rl.RayWhite)
		drawMenu()
	case Playing:
		rl.ClearBackground(rl.SkyBlue)
		camera := rl.Camera3D{
			Position: rl.Vector3{
				X: car.position.X - 5*float32(math.Cos(float64(car.yaw))),
				Y: car.position.Y + 2,
				Z: car.position.Z - 5*float32(math.Sin(float64(car.yaw))),
			},
			Target:     car.position,
			Up:         rl.Vector3{Y: 1},
			Fovy:       45,
			Projection: rl.CameraPerspective,
		}
		rl.BeginMode3D(camera)
		drawWorld()
		drawCar()
		rl.EndMode3D()
	}
}

func handleMenuInput() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePos := rl.GetMousePosition()
		buttonX := float32(rl.GetScreenWidth()-200) / 2
		buttonY := float32(rl.GetScreenHeight()-50) / 2
		if mousePos.X >= buttonX && mousePos.X <= buttonX+200 &&
			mousePos.Y >= buttonY && mousePos.Y <= buttonY+50 {
			currentState = Playing
		}
	}
}

func drawMenu() {
	buttonX := (rl.GetScreenWidth() - 200) / 2
	buttonY := (rl.GetScreenHeight() - 50) / 2
	rl.DrawRectangle(int32(buttonX), int32(buttonY), 200, 50, rl.Gray)
	rl.DrawText("Play", int32(buttonX+50), int32(buttonY+10), 30, rl.Black)
}

func updatePhysics() {
	// Simplified ramp physics (assuming ramps are handled elsewhere or removed for simplicity)
	if car.position.Y <= 0 {
		car.position.Y = 0
		car.pitch = 0
		car.grounded = true
	} else {
		car.velocity.Y -= 0.5 * rl.GetFrameTime()
		car.position.Y += car.velocity.Y * rl.GetFrameTime()
		car.grounded = false
	}
}

func checkCollisions() {
	// Define the car's bounding box (assuming car dimensions are 1x0.5x2)
	carHalfSize := rl.Vector3{X: 0.5, Y: 0.25, Z: 1}
	carMin := rl.Vector3Subtract(car.position, carHalfSize)
	carMax := rl.Vector3Add(car.position, carHalfSize)

	// Check collision with barriers (existing code)
	for _, barrier := range barrierModels {
		// Extract position from the barrier's Transform matrix
		barrierPos := rl.Vector3{
			X: barrier.Transform.M12, // Translation X
			Y: barrier.Transform.M13, // Translation Y
			Z: barrier.Transform.M14, // Translation Z
		}
		barrierHalfSize := rl.Vector3{X: 0.5, Y: 0.5, Z: 0.5} // Barrier size: 1x1x1
		barrierMin := rl.Vector3Subtract(barrierPos, barrierHalfSize)
		barrierMax := rl.Vector3Add(barrierPos, barrierHalfSize)

		// AABB collision check for barriers
		if carMin.X < barrierMax.X && carMax.X > barrierMin.X &&
			carMin.Y < barrierMax.Y && carMax.Y > barrierMin.Y &&
			carMin.Z < barrierMax.Z && carMax.Z > barrierMin.Z {
			car.position = prevCarPosition // Revert to previous position
			return                         // Exit early on collision
		}
	}

	// Check collision with the building
	buildingPos := rl.Vector3{X: 50, Y: 25, Z: 50}    // Building center
	buildingHalfSize := rl.Vector3{X: 5, Y: 25, Z: 5} // Half of 10x50x10
	buildingMin := rl.Vector3Subtract(buildingPos, buildingHalfSize)
	buildingMax := rl.Vector3Add(buildingPos, buildingHalfSize)

	// AABB collision check for the building
	if carMin.X < buildingMax.X && carMax.X > buildingMin.X &&
		carMin.Y < buildingMax.Y && carMax.Y > buildingMin.Y &&
		carMin.Z < buildingMax.Z && carMax.Z > buildingMin.Z {
		car.position = prevCarPosition // Revert to previous position
	}
}
