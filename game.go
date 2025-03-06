package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	Menu GameState = iota
	Playing
)

var currentState GameState

// In-game settings overlay state.
var showSettingsOverlay bool = false

// Toggle for displaying the FPS counter.
var showFPSCounter bool = true

func initGame() {
	currentState = Menu
	initCar()
	initWorld()
	// Ensure settings overlay is off when starting
	showSettingsOverlay = false
}

func updateGame() {
	switch currentState {
	case Menu:
		// Main menu: Only a "Play" button.
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mousePos := rl.GetMousePosition()
			screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
			playX, playY := float32(screenW/2-100), float32(screenH/2-25)
			if mousePos.X >= playX && mousePos.X <= playX+200 &&
				mousePos.Y >= playY && mousePos.Y <= playY+50 {
				currentState = Playing
			}
		}
	case Playing:
		// In-game, first update settings overlay UI:
		mousePos := rl.GetMousePosition()
		// Gear icon area (top left 40x40)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if mousePos.X >= 10 && mousePos.X <= 50 &&
				mousePos.Y >= 10 && mousePos.Y <= 50 {
				// Toggle settings overlay.
				showSettingsOverlay = !showSettingsOverlay
			}
		}

		// If settings overlay is open, check its buttons.
		if showSettingsOverlay {
			// We'll draw a panel in the center with two buttons.
			// Button for toggling FPS counter and button for returning to main menu.
			screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
			// Define panel rectangle (centered, 300x200)
			panelX, panelY := float32(screenW/2-150), float32(screenH/2-100)
			// Button positions relative to panel.
			toggleFPSX, toggleFPSY := panelX+50, panelY+50
			returnX, returnY := panelX+50, panelY+120

			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				// Toggle FPS button.
				if mousePos.X >= toggleFPSX && mousePos.X <= toggleFPSX+200 &&
					mousePos.Y >= toggleFPSY && mousePos.Y <= toggleFPSY+40 {
					showFPSCounter = !showFPSCounter
				}
				// Return to Main Menu button.
				if mousePos.X >= returnX && mousePos.X <= returnX+200 &&
					mousePos.Y >= returnY && mousePos.Y <= returnY+40 {
					currentState = Menu
					// Reset in-game state.
					showSettingsOverlay = false
				}
			}
		} else {
			// Only update car if settings overlay is not open.
			updateCar()
		}
	}
}

func drawGame() {
	switch currentState {
	case Menu:
		rl.ClearBackground(rl.RayWhite)
		screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
		playX, playY := (screenW-200)/2, (screenH-50)/2
		rl.DrawRectangle(int32(playX), int32(playY), 200, 50, rl.Gray)
		rl.DrawText("Play", int32(playX+70), int32(playY+10), 30, rl.Black)
	case Playing:
		rl.ClearBackground(rl.SkyBlue)
		// Set up camera.
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

		// Draw gear icon (a simple square with a gear symbol)
		rl.DrawRectangle(10, 10, 40, 40, rl.Gray)
		rl.DrawText("⚙", 20, 10, 32, rl.Black) // You can change the symbol if needed.

		// Draw FPS counter in top right if enabled.
		if showFPSCounter {
			fpsText := fmt.Sprintf("FPS: %d", rl.GetFPS())
			screenW := rl.GetScreenWidth()
			rl.DrawText(fpsText, int32(screenW)-100, 10, 20, rl.Black)
		}

		// If settings overlay is open, draw it.
		if showSettingsOverlay {
			screenW, screenH := rl.GetScreenWidth(), rl.GetScreenHeight()
			panelX, panelY := (screenW-300)/2, (screenH-200)/2
			rl.DrawRectangle(int32(panelX), int32(panelY), 300, 200, rl.Fade(rl.LightGray, 0.9))
			rl.DrawText("Settings", int32(panelX+100), int32(panelY+20), 30, rl.Black)
			// Toggle FPS button.
			toggleFPSX, toggleFPSY := panelX+50, panelY+50
			rl.DrawRectangle(int32(toggleFPSX), int32(toggleFPSY), 200, 40, rl.Gray)
			toggleText := "FPS: OFF"
			if showFPSCounter {
				toggleText = "FPS: ON"
			}
			rl.DrawText(toggleText, int32(toggleFPSX+60), int32(toggleFPSY+10), 20, rl.Black)
			// Return to Main Menu button.
			returnX, returnY := panelX+50, panelY+120
			rl.DrawRectangle(int32(returnX), int32(returnY), 200, 40, rl.Gray)
			rl.DrawText("Return to Main Menu", int32(returnX+10), int32(returnY+10), 20, rl.Black)
		}
	}
}
