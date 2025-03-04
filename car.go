package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Car struct {
	position rl.Vector3
	yaw      float32
	pitch    float32
	speed    float32
	steering float32
	velocity rl.Vector3
	grounded bool
	model    rl.Model
}

var car Car // Global car variable

func initCar() {
	car = Car{
		position: rl.Vector3{X: 0, Y: 0, Z: 0},
		yaw:      0,
		pitch:    0,
		speed:    0,
		steering: 0,
		grounded: true,
		model:    rl.LoadModelFromMesh(rl.GenMeshCube(1, 0.5, 2)),
	}
}

// checkCollisions returns true if the car collides with any collider.
func checkCollisions(pos rl.Vector3) bool {
	for _, c := range colliders {
		dx := pos.X - c.Position.X
		dz := pos.Z - c.Position.Z
		distance := math.Sqrt(float64(dx*dx + dz*dz))
		// Assuming car collision radius is 1. Adjust as needed.
		if float32(distance) < (1 + c.Radius) {
			return true
		}
	}
	return false
}

func updateCar() {
	oldPos := car.position // save previous position

	if rl.IsKeyDown(rl.KeyUp) {
		car.speed += 5 * rl.GetFrameTime()
		if car.speed > 10 {
			car.speed = 10
		}
	} else if rl.IsKeyDown(rl.KeyDown) {
		car.speed -= 5 * rl.GetFrameTime()
		if car.speed < -5 {
			car.speed = -5
		}
	} else {
		if car.speed > 0 {
			car.speed -= 2 * rl.GetFrameTime()
			if car.speed < 0 {
				car.speed = 0
			}
		} else if car.speed < 0 {
			car.speed += 2 * rl.GetFrameTime()
			if car.speed > 0 {
				car.speed = 0
			}
		}
	}

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
	car.position.X += forward.X * car.speed * rl.GetFrameTime()
	car.position.Z += forward.Z * car.speed * rl.GetFrameTime()

	// Check collisions. If collision, revert position and zero speed.
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
