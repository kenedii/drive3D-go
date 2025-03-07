package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Car struct {
	position rl.Vector3
	yaw      float32
	pitch    float32
	speed    float32 // in m/s
	steering float32
	velocity rl.Vector3
	grounded bool
	model    rl.Model
}

var car Car

func initCar() {
	// Spawn on road: center of chunk (0,0) is at (CHUNK_SIZE/2, 0, CHUNK_SIZE/2)
	car = Car{
		position: rl.Vector3{X: CHUNK_SIZE / 2, Y: 0, Z: CHUNK_SIZE / 2},
		yaw:      0,
		pitch:    0,
		speed:    0,
		steering: 0,
		grounded: true,
		model:    rl.LoadModelFromMesh(rl.GenMeshCube(1, 0.5, 2)),
	}
}

func updateCar() {
	oldPos := car.position

	// Determine terrain multiplier and max speed based on road type.
	var terrainMultiplier, maxSpeed float32
	chunk := chunks[getChunkCoord(car.position)]
	if chunk != nil {
		switch chunk.RoadType {
		case RoadIce:
			terrainMultiplier = 1.5
			maxSpeed = 50.0 // ≈180 km/h
		case RoadDirt:
			terrainMultiplier = 0.85
			maxSpeed = 25.0 // ≈90 km/h
		default:
			terrainMultiplier = 1.0
			maxSpeed = 38.9 // ≈140 km/h
		}
	} else {
		terrainMultiplier = 1.0
		maxSpeed = 38.9 // Default road speed if no chunk data
	}

	// Base acceleration
	baseAccel := float32(5.0) // m/s^2
	accel := baseAccel * terrainMultiplier

	// Define a deceleration factor for overspeed.
	overspeedDecelFactor := float32(2.0) // Stronger deceleration

	dt := rl.GetFrameTime()

	// Process acceleration input.
	if rl.IsKeyDown(rl.KeyUp) {
		car.speed += accel * dt
		// Gradually slow down if exceeding max speed
		if car.speed > maxSpeed {
			car.speed -= overspeedDecelFactor * (car.speed - maxSpeed) * dt
		}
	} else if rl.IsKeyDown(rl.KeyDown) {
		car.speed -= accel * dt
		if car.speed < -maxSpeed/2 {
			car.speed += overspeedDecelFactor * (-maxSpeed/2 - car.speed) * dt
		}
	} else {
		// Apply friction to gradually slow down momentum
		var friction float32
		if chunk != nil {
			if chunk.RoadType == RoadNormal {
				friction = 0.995
			} else if chunk.RoadType == RoadIce {
				friction = 0.999 // Very slow deceleration on ice
			} else {
				friction = 0.98 // Stronger deceleration on dirt
			}
		} else {
			friction = 0.995
		}

		car.speed *= friction
		if math.Abs(float64(car.speed)) < 0.1 {
			car.speed = 0
		}
	}

	// Steering
	if rl.IsKeyDown(rl.KeyLeft) {
		car.steering -= 2 * dt
		if car.steering < -1 {
			car.steering = -1
		}
	} else if rl.IsKeyDown(rl.KeyRight) {
		car.steering += 2 * dt
		if car.steering > 1 {
			car.steering = 1
		}
	} else {
		car.steering *= 0.9
	}

	car.yaw += car.steering * dt

	// Compute forward direction
	forward := rl.Vector3{
		X: float32(math.Cos(float64(car.yaw))) * float32(math.Cos(float64(car.pitch))),
		Y: float32(math.Sin(float64(car.pitch))),
		Z: float32(math.Sin(float64(car.yaw))) * float32(math.Cos(float64(car.pitch))),
	}

	// Update position
	car.position.X += forward.X * car.speed * dt
	car.position.Z += forward.Z * car.speed * dt

	// Collision check
	if checkCollisions(car.position) {
		car.position = oldPos
		car.speed = 0
	}

	// Grounded check
	if car.grounded {
		car.velocity.Y = 0
	}
}

func drawCar() {
	trans := rl.MatrixTranslate(car.position.X, car.position.Y, car.position.Z)
	rotY := rl.MatrixRotateY(car.yaw)
	rotX := rl.MatrixRotateX(car.pitch)
	transform := rl.MatrixMultiply(rotX, rotY)
	transform = rl.MatrixMultiply(transform, trans)
	car.model.Transform = transform
	rl.DrawModel(car.model, rl.Vector3{}, 1, rl.Red)
}
