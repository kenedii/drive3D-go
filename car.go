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

	// Determine terrain multiplier based on current chunk's road type.
	// Normal roads: 1.0, Dirt roads: 0.85, Ice roads: 1.5.
	terrainMultiplier := float32(1.0)
	chunk := chunks[getChunkCoord(car.position)]
	if chunk != nil {
		switch chunk.RoadType {
		case RoadIce:
			terrainMultiplier = 1.5
		case RoadDirt:
			terrainMultiplier = 0.85
		default:
			terrainMultiplier = 1.0
		}
	}

	// Define base acceleration and maximum speed for normal roads.
	// Increase baseMaxSpeed from 10 m/s to 15 m/s.
	baseAccel := float32(5.0)     // m/s^2
	baseMaxSpeed := float32(15.0) // m/s (≈54 km/h on normal roads)
	baseDecel := float32(2.0)     // m/s^2

	// Apply terrain multiplier to acceleration and maximum speed.
	accel := baseAccel * terrainMultiplier
	maxSpeed := baseMaxSpeed * terrainMultiplier

	// Update speed based on input.
	if rl.IsKeyDown(rl.KeyUp) {
		car.speed += accel * rl.GetFrameTime()
		if car.speed > maxSpeed {
			car.speed = maxSpeed
		}
	} else if rl.IsKeyDown(rl.KeyDown) {
		car.speed -= accel * rl.GetFrameTime()
		if car.speed < -maxSpeed/2 { // allow slower reverse
			car.speed = -maxSpeed / 2
		}
	} else {
		// Deceleration.
		if car.speed > 0 {
			car.speed -= baseDecel * rl.GetFrameTime()
			if car.speed < 0 {
				car.speed = 0
			}
		} else if car.speed < 0 {
			car.speed += baseDecel * rl.GetFrameTime()
			if car.speed > 0 {
				car.speed = 0
			}
		}
	}

	// Update steering (steering multiplier remains at 1.0).
	if rl.IsKeyDown(rl.KeyLeft) {
		car.steering -= 2 * rl.GetFrameTime()
		if car.steering < -1 {
			car.steering = -1
		}
	} else if rl.IsKeyDown(rl.KeyRight) {
		car.steering += 2 * rl.GetFrameTime()
		if car.steering > 1 {
			car.steering = 1
		}
	} else {
		car.steering *= 0.9
	}

	car.yaw += car.steering * rl.GetFrameTime()

	forward := rl.Vector3{
		X: float32(math.Cos(float64(car.yaw))) * float32(math.Cos(float64(car.pitch))),
		Y: float32(math.Sin(float64(car.pitch))),
		Z: float32(math.Sin(float64(car.yaw))) * float32(math.Cos(float64(car.pitch))),
	}

	// Update position based on speed and direction.
	car.position.X += forward.X * car.speed * rl.GetFrameTime()
	car.position.Z += forward.Z * car.speed * rl.GetFrameTime()

	if checkCollisions(car.position) {
		car.position = oldPos
		car.speed = 0
	}

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
