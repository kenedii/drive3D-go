package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 600, "3D Racing Game")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	initGame()

	for !rl.WindowShouldClose() {
		updateGame()
		rl.BeginDrawing()
		drawGame()
		rl.EndDrawing()
	}
}
