package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed RNG for procedural generation
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
