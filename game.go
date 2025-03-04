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
	currentState GameState
)

func initGame() {
	currentState = Menu
	initCar()   // From car.go
	initWorld() // From world.go
}

func updateGame() {
	switch currentState {
	case Menu:
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePos := rl.GetMousePosition()
			buttonX := float32(rl.GetScreenWidth()-200) / 2
			buttonY := float32(rl.GetScreenHeight()-50) / 2
			if mousePos.X >= buttonX && mousePos.X <= buttonX+200 &&
				mousePos.Y >= buttonY && mousePos.Y <= buttonY+50 {
				currentState = Playing
			}
		}
	case Playing:
		updateCar() // From car.go
		// updateWorld is now called within drawWorld(), so additional call here is optional
	}
}

func drawGame() {
	switch currentState {
	case Menu:
		rl.ClearBackground(rl.RayWhite)
		buttonX := (rl.GetScreenWidth() - 200) / 2
		buttonY := (rl.GetScreenHeight() - 50) / 2
		rl.DrawRectangle(int32(buttonX), int32(buttonY), 200, 50, rl.Gray)
		rl.DrawText("Play", int32(buttonX+50), int32(buttonY+10), 30, rl.Black)
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
		drawWorld() // From world.go
		drawCar()   // From car.go
		rl.EndMode3D()
	}
}
