package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	groundModel     rl.Model
	roadModel       rl.Model
	barrierModels   []rl.Model
	skyscraperModel rl.Model
)

func initWorld() {
	// Ground: 200x200 units
	groundModel = rl.LoadModelFromMesh(rl.GenMeshPlane(200, 200, 1, 1))

	// Road: 100x5 units
	roadModel = rl.LoadModelFromMesh(rl.GenMeshPlane(100, 5, 1, 1))
	roadModel.Transform = rl.MatrixTranslate(0, 0.01, 0)

	// Skyscraper: 10x50x10 units at (50, 25, 50)
	skyscraperModel = rl.LoadModelFromMesh(rl.GenMeshCube(10, 50, 10))
	skyscraperModel.Transform = rl.MatrixTranslate(50, 25, 50)

	// Barriers: 1x1x1 cubes along edges
	barrierModels = make([]rl.Model, 0)
	barrierSize := float32(1)
	landscapeSize := float32(200)

	// Barriers along X-axis (left and right edges)
	for x := -landscapeSize / 2; x <= landscapeSize/2; x += barrierSize {
		barrierLeft := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
		barrierLeft.Transform = rl.MatrixTranslate(x, 0.5, -landscapeSize/2)
		barrierModels = append(barrierModels, barrierLeft)

		barrierRight := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
		barrierRight.Transform = rl.MatrixTranslate(x, 0.5, landscapeSize/2)
		barrierModels = append(barrierModels, barrierRight)
	}

	// Barriers along Z-axis (top and bottom edges)
	for z := -landscapeSize/2 + barrierSize; z < landscapeSize/2; z += barrierSize {
		barrierTop := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
		barrierTop.Transform = rl.MatrixTranslate(-landscapeSize/2, 0.5, z)
		barrierModels = append(barrierModels, barrierTop)

		barrierBottom := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
		barrierBottom.Transform = rl.MatrixTranslate(landscapeSize/2, 0.5, z)
		barrierModels = append(barrierModels, barrierBottom)
	}
}

func drawWorld() {
	rl.DrawModel(groundModel, rl.Vector3{Y: -0.01}, 1, rl.Green)
	rl.DrawModel(roadModel, rl.Vector3{}, 1, rl.DarkGray)
	rl.DrawModel(skyscraperModel, rl.Vector3{}, 1, rl.Blue) // Draw skyscraper
	for _, barrier := range barrierModels {
		rl.DrawModel(barrier, rl.Vector3{}, 1, rl.Red)
	}
}
