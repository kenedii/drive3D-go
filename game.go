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

var currentState GameState

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
		updateCar()
		updatePhysics()
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
	// Ramp 1: along x, x=5 to 15, z=-1 to 1, y=0 to 5
	if car.position.X >= 5 && car.position.X <= 15 && car.position.Z >= -1 && car.position.Z <= 1 {
		car.position.Y = (car.position.X - 5) / 10 * 5
		car.pitch = 26.565 * rl.Deg2rad // Convert degrees to radians
		car.grounded = true
	} else if car.position.Z >= 5 && car.position.Z <= 15 && car.position.X >= -1 && car.position.X <= 1 {
		// Ramp 2: along z, z=5 to 15, x=-1 to 1, y=0 to 5
		car.position.Y = (car.position.Z - 5) / 10 * 5
		car.pitch = 26.565 * rl.Deg2rad // Convert degrees to radians
		car.grounded = true
	} else if car.position.Y <= 0 {
		car.position.Y = 0
		car.pitch = 0
		car.grounded = true
	} else {
		// In air, apply gravity
		car.velocity.Y -= 0.5 * rl.GetFrameTime()
		car.position.Y += car.velocity.Y * rl.GetFrameTime()
		car.grounded = false
	}
}
