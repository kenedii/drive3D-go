package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	Menu GameState = iota
	Playing
	Settings
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
		// Two buttons: "Play" and "Settings"
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePos := rl.GetMousePosition()
			screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
			playX, playY := float32(screenW/2-100), float32(screenH/2-60)
			setX, setY := float32(screenW/2-100), float32(screenH/2+20)
			if mousePos.X >= playX && mousePos.X <= playX+200 &&
				mousePos.Y >= playY && mousePos.Y <= playY+50 {
				currentState = Playing
			} else if mousePos.X >= setX && mousePos.X <= setX+200 &&
				mousePos.Y >= setY && mousePos.Y <= setY+50 {
				currentState = Settings
			}
		}
	case Playing:
		updateCar()
		// (In Playing state, you might add pause handling later.)
	case Settings:
		// In Settings, if "Return" is clicked, go back to Menu.
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePos := rl.GetMousePosition()
			screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
			returnX, returnY := float32(screenW/2-100), float32(screenH/2+20)
			if mousePos.X >= returnX && mousePos.X <= returnX+200 &&
				mousePos.Y >= returnY && mousePos.Y <= returnY+50 {
				currentState = Menu
			}
		}
	}
}

func drawGame() {
	switch currentState {
	case Menu:
		rl.ClearBackground(rl.RayWhite)
		screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
		playX, playY := (screenW-200)/2, (screenH-50)/2-60
		setX, setY := (screenW-200)/2, (screenH-50)/2+20
		rl.DrawRectangle(int32(playX), int32(playY), 200, 50, rl.Gray)
		rl.DrawText("Play", int32(playX+70), int32(playY+10), 30, rl.Black)
		rl.DrawRectangle(int32(setX), int32(setY), 200, 50, rl.Gray)
		rl.DrawText("Settings", int32(setX+40), int32(setY+10), 30, rl.Black)
	case Playing:
		rl.ClearBackground(rl.SkyBlue)
		camera := rl.Camera3D{
			Position: rl.Vector3{
				X: car.position.X - 5,
				Y: car.position.Y + 2,
				Z: car.position.Z - 5,
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
	case Settings:
		rl.ClearBackground(rl.RayWhite)
		screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
		rl.DrawText("Settings", int32(screenW/2-60), int32(screenH/2-100), 40, rl.Black)
		// Display FPS.
		fpsText := fmt.Sprintf("FPS: %d", rl.GetFPS())
		rl.DrawText(fpsText, int32(screenW/2-40), int32(screenH/2-40), 30, rl.Black)
		// Return button.
		returnX, returnY := (screenW-200)/2, (screenH-50)/2+20
		rl.DrawRectangle(int32(returnX), int32(returnY), 200, 50, rl.Gray)
		rl.DrawText("Return", int32(returnX+50), int32(returnY+10), 30, rl.Black)
	}
}
