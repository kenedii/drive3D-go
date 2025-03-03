package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	groundModel   rl.Model
	roadModel     rl.Model
	barrierModels []rl.Model
)

func initWorld() {
	// Expanded ground plane: 200x200 units
	groundModel = rl.LoadModelFromMesh(rl.GenMeshPlane(200, 200, 1, 1))

	// Road: 100 units long, 5 units wide, positioned along X-axis
	roadMesh := rl.GenMeshPlane(100, 5, 1, 1)
	roadModel = rl.LoadModelFromMesh(roadMesh)
	roadModel.Transform = rl.MatrixTranslate(0, 0.01, 0) // Slightly above ground to avoid z-fighting

	// Barriers: cubes along the edges
	barrierModels = make([]rl.Model, 0)
	barrierSize := float32(1)
	barrierHeight := float32(1)
	landscapeSize := float32(200)

	// Barriers along X-axis (left and right edges)
	for x := -landscapeSize / 2; x <= landscapeSize/2; x += barrierSize {
		// Left barrier
		barrierLeft := rl.LoadModelFromMesh(rl.GenMeshCube(barrierSize, barrierHeight, barrierSize))
		barrierLeft.Transform = rl.MatrixTranslate(x, barrierHeight/2, -landscapeSize/2)
		barrierModels = append(barrierModels, barrierLeft)

		// Right barrier
		barrierRight := rl.LoadModelFromMesh(rl.GenMeshCube(barrierSize, barrierHeight, barrierSize))
		barrierRight.Transform = rl.MatrixTranslate(x, barrierHeight/2, landscapeSize/2)
		barrierModels = append(barrierModels, barrierRight)
	}

	// Barriers along Z-axis (top and bottom edges)
	for z := -landscapeSize/2 + barrierSize; z < landscapeSize/2; z += barrierSize {
		// Top barrier
		barrierTop := rl.LoadModelFromMesh(rl.GenMeshCube(barrierSize, barrierHeight, barrierSize))
		barrierTop.Transform = rl.MatrixTranslate(-landscapeSize/2, barrierHeight/2, z)
		barrierModels = append(barrierModels, barrierTop)

		// Bottom barrier
		barrierBottom := rl.LoadModelFromMesh(rl.GenMeshCube(barrierSize, barrierHeight, barrierSize))
		barrierBottom.Transform = rl.MatrixTranslate(landscapeSize/2, barrierHeight/2, z)
		barrierModels = append(barrierModels, barrierBottom)
	}
}

func drawWorld() {
	// Draw ground
	rl.DrawModel(groundModel, rl.Vector3{Y: -0.01}, 1, rl.Green)

	// Draw road
	rl.DrawModel(roadModel, rl.Vector3{}, 1, rl.DarkGray)

	// Draw barriers
	for _, barrier := range barrierModels {
		rl.DrawModel(barrier, rl.Vector3{}, 1, rl.Red)
	}
}
