package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Use float32 for CHUNK_SIZE.
const CHUNK_SIZE float32 = 50.0

var LightBlue = rl.Color{R: 173, G: 216, B: 230, A: 255} // RGB for light blue

// Chunk types.
const (
	Highway    = 0
	City       = 1
	Commercial = 2
	Desert     = 3
	Forest     = 4
	Snow       = 5
)

// Road types.
const (
	RoadNormal = 0
	RoadDirt   = 1
	RoadIce    = 2
)

// Coord represents a chunk's position.
type Coord struct {
	X, Y int
}

// Chunk holds the chunk type, road type, and models for rendering.
type Chunk struct {
	Type     int
	RoadType int
	Coord    Coord
	Models   []rl.Model
}

var (
	chunks          map[Coord]*Chunk
	lastPlayerChunk Coord
	// Allowed neighbors for each chunk type.
	allowedNeighbors = [][]int{
		{Highway, City, Commercial, Desert, Forest, Snow}, // Highway
		{Highway, Commercial},                             // City
		{Highway, City},                                   // Commercial
		{Highway, Desert},                                 // Desert
		{Highway, Forest},                                 // Forest
		{Highway, Snow},                                   // Snow
	}
	// Colors for ground based on chunk type.
	typeColors = []rl.Color{
		rl.Gray,     // Highway
		rl.DarkGray, // City
		rl.Gray,     // Commercial
		rl.Yellow,   // Desert
		rl.Green,    // Forest
		rl.White,    // Snow
	}
)

// collisionBoxes holds bounding boxes for objects.
var collisionBoxes []rl.BoundingBox

// getChunkCoord converts a world position to chunk coordinates.
func getChunkCoord(pos rl.Vector3) Coord {
	i := int(math.Floor(float64(pos.X / CHUNK_SIZE)))
	j := int(math.Floor(float64(pos.Z / CHUNK_SIZE)))
	return Coord{i, j}
}

// checkCollisions returns true if the given position collides with any bounding box.
func checkCollisions(pos rl.Vector3) bool {
	// Treat the car as a circle with radius 1 (XZ plane).
	for _, box := range collisionBoxes {
		centerX := (box.Min.X + box.Max.X) / 2
		centerZ := (box.Min.Z + box.Max.Z) / 2
		dx := pos.X - centerX
		dz := pos.Z - centerZ
		distance := math.Sqrt(float64(dx*dx + dz*dz))
		boxRadius := (box.Max.X - box.Min.X) / 2
		if distance < float64(1+boxRadius) {
			return true
		}
	}
	return false
}

// setAlbedoColor updates the model's material color.
func setAlbedoColor(model *rl.Model, color rl.Color) {
	mat := *model.Materials // copy the material
	mat.Maps.Color = color
	model.Materials = &mat
}

// isPositionOnRoad returns true if (x,z) (relative to the chunk origin) lies on the road.
// In our design the road is a plus shape through the center.
func isPositionOnRoad(x, z float32) bool {
	roadWidth := float32(5.0)
	center := CHUNK_SIZE / 2
	if math.Abs(float64(x-center)) <= float64(roadWidth) || math.Abs(float64(z-center)) <= float64(roadWidth) {
		return true
	}
	return false
}

// determineChunkType returns a chunk type based on neighbors.
func determineChunkType(i, j int) int {
	neighbors := []Coord{{i - 1, j}, {i + 1, j}, {i, j - 1}, {i, j + 1}}
	var types []int
	for _, n := range neighbors {
		if chunk, exists := chunks[n]; exists {
			types = append(types, chunk.Type)
		}
	}
	if len(types) == 0 {
		return rand.Intn(6)
	}
	allowed := make(map[int]bool)
	for _, t := range allowedNeighbors[types[0]] {
		allowed[t] = true
	}
	for _, t := range types[1:] {
		newAllowed := make(map[int]bool)
		for _, a := range allowedNeighbors[t] {
			if allowed[a] {
				newAllowed[a] = true
			}
		}
		allowed = newAllowed
	}
	var allowedTypes []int
	for t := range allowed {
		allowedTypes = append(allowedTypes, t)
	}
	if len(allowedTypes) == 0 {
		return Highway
	}
	return allowedTypes[rand.Intn(len(allowedTypes))]
}

// generateChunk creates a chunk at grid coordinate (i,j) with ground, road, and objects.
func generateChunk(i, j int) {
	coord := Coord{i, j}
	if _, exists := chunks[coord]; exists {
		return
	}
	chunkType := determineChunkType(i, j)
	// Determine road type based on chunk type.
	var roadType int
	switch chunkType {
	case Forest:
		roadType = RoadDirt
	case Snow:
		roadType = RoadIce
	case Highway:
		if rand.Float32() < 0.3 {
			roadType = RoadDirt
		} else {
			roadType = RoadNormal
		}
	default:
		roadType = RoadNormal
	}
	chunk := &Chunk{Type: chunkType, RoadType: roadType, Coord: coord, Models: []rl.Model{}}
	posX := float32(i) * CHUNK_SIZE
	posZ := float32(j) * CHUNK_SIZE

	// Ground.
	groundMesh := rl.GenMeshPlane(CHUNK_SIZE, CHUNK_SIZE, 1, 1)
	groundModel := rl.LoadModelFromMesh(groundMesh)
	groundTransform := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0, posZ+CHUNK_SIZE/2)
	groundModel.Transform = groundTransform
	setAlbedoColor(&groundModel, typeColors[chunkType])
	chunk.Models = append(chunk.Models, groundModel)

	// Choose road color based on road type.
	var roadColor rl.Color
	switch roadType {
	case RoadNormal:
		roadColor = rl.DarkGray
	case RoadDirt:
		roadColor = rl.NewColor(139, 69, 19, 255) // brownish
	case RoadIce:
		roadColor = LightBlue
	}

	// Main "+" road.
	roadWidth := float32(5.0)
	roadMeshH := rl.GenMeshPlane(CHUNK_SIZE, roadWidth, 1, 1)
	roadModelH := rl.LoadModelFromMesh(roadMeshH)
	roadTransformH := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+CHUNK_SIZE/2)
	roadModelH.Transform = roadTransformH
	setAlbedoColor(&roadModelH, roadColor)
	chunk.Models = append(chunk.Models, roadModelH)

	roadMeshV := rl.GenMeshPlane(roadWidth, CHUNK_SIZE, 1, 1)
	roadModelV := rl.LoadModelFromMesh(roadMeshV)
	roadTransformV := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+CHUNK_SIZE/2)
	roadModelV.Transform = roadTransformV
	setAlbedoColor(&roadModelV, roadColor)
	chunk.Models = append(chunk.Models, roadModelV)

	// For central chunk (0,0), only ground and main road.
	if i == 0 && j == 0 {
		chunks[coord] = chunk
		return
	}

	// For City/Commercial chunks, add backroads.
	if chunkType == City || chunkType == Commercial {
		roadModelTop := rl.LoadModelFromMesh(roadMeshH)
		roadTopTransform := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+CHUNK_SIZE-roadWidth)
		roadModelTop.Transform = roadTopTransform
		setAlbedoColor(&roadModelTop, roadColor)
		chunk.Models = append(chunk.Models, roadModelTop)

		roadModelBottom := rl.LoadModelFromMesh(roadMeshH)
		roadBottomTransform := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+roadWidth)
		roadModelBottom.Transform = roadBottomTransform
		setAlbedoColor(&roadModelBottom, roadColor)
		chunk.Models = append(chunk.Models, roadModelBottom)

		roadModelLeft := rl.LoadModelFromMesh(roadMeshV)
		roadLeftTransform := rl.MatrixTranslate(posX+roadWidth, 0.01, posZ+CHUNK_SIZE/2)
		roadModelLeft.Transform = roadLeftTransform
		setAlbedoColor(&roadModelLeft, roadColor)
		chunk.Models = append(chunk.Models, roadModelLeft)

		roadModelRight := rl.LoadModelFromMesh(roadMeshV)
		roadRightTransform := rl.MatrixTranslate(posX+CHUNK_SIZE-roadWidth, 0.01, posZ+CHUNK_SIZE/2)
		roadModelRight.Transform = roadRightTransform
		setAlbedoColor(&roadModelRight, roadColor)
		chunk.Models = append(chunk.Models, roadModelRight)
	}

	// Seeded randomness for object placement.
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%d,%d", i, j)))
	seed := h.Sum64()
	chunkRand := rand.New(rand.NewSource(int64(seed)))

	// Spawn objects based on chunk type.
	switch chunkType {
	case City:
		for k := 0; k < 5; k++ {
			bx := posX + chunkRand.Float32()*CHUNK_SIZE
			bz := posZ + chunkRand.Float32()*CHUNK_SIZE
			if !isPositionOnRoad(bx-posX, bz-posZ) {
				buildingMesh := rl.GenMeshCube(10, 50, 10)
				buildingModel := rl.LoadModelFromMesh(buildingMesh)
				buildingTransform := rl.MatrixTranslate(bx, 25, bz)
				buildingModel.Transform = buildingTransform
				setAlbedoColor(&buildingModel, rl.Blue)
				chunk.Models = append(chunk.Models, buildingModel)
				box := rl.BoundingBox{
					Min: rl.Vector3{X: bx - 5, Y: 0, Z: bz - 5},
					Max: rl.Vector3{X: bx + 5, Y: 50, Z: bz + 5},
				}
				collisionBoxes = append(collisionBoxes, box)
			}
		}
	case Commercial:
		for k := 0; k < 3; k++ {
			sx := posX + chunkRand.Float32()*CHUNK_SIZE
			sz := posZ + chunkRand.Float32()*CHUNK_SIZE
			if !isPositionOnRoad(sx-posX, sz-posZ) {
				storeMesh := rl.GenMeshCube(15, 10, 15)
				storeModel := rl.LoadModelFromMesh(storeMesh)
				storeTransform := rl.MatrixTranslate(sx, 5, sz)
				storeModel.Transform = storeTransform
				setAlbedoColor(&storeModel, rl.Purple)
				chunk.Models = append(chunk.Models, storeModel)
				box := rl.BoundingBox{
					Min: rl.Vector3{X: sx - 7.5, Y: 0, Z: sz - 7.5},
					Max: rl.Vector3{X: sx + 7.5, Y: 10, Z: sz + 7.5},
				}
				collisionBoxes = append(collisionBoxes, box)
			}
		}
	case Desert:
		for k := 0; k < 10; k++ {
			cx := posX + chunkRand.Float32()*CHUNK_SIZE
			cz := posZ + chunkRand.Float32()*CHUNK_SIZE
			if !isPositionOnRoad(cx-posX, cz-posZ) {
				cactusMesh := rl.GenMeshCube(1, 5, 1)
				cactusModel := rl.LoadModelFromMesh(cactusMesh)
				cactusTransform := rl.MatrixTranslate(cx, 2.5, cz)
				cactusModel.Transform = cactusTransform
				setAlbedoColor(&cactusModel, rl.Green)
				chunk.Models = append(chunk.Models, cactusModel)
				box := rl.BoundingBox{
					Min: rl.Vector3{X: cx - 0.5, Y: 0, Z: cz - 0.5},
					Max: rl.Vector3{X: cx + 0.5, Y: 5, Z: cz + 0.5},
				}
				collisionBoxes = append(collisionBoxes, box)
			}
		}
	case Forest:
		for k := 0; k < 20; k++ {
			tx := posX + chunkRand.Float32()*CHUNK_SIZE
			tz := posZ + chunkRand.Float32()*CHUNK_SIZE
			if !isPositionOnRoad(tx-posX, tz-posZ) {
				treeMesh := rl.GenMeshCube(2, 10, 2)
				treeModel := rl.LoadModelFromMesh(treeMesh)
				treeTransform := rl.MatrixTranslate(tx, 5, tz)
				treeModel.Transform = treeTransform
				setAlbedoColor(&treeModel, rl.DarkGreen)
				chunk.Models = append(chunk.Models, treeModel)
				box := rl.BoundingBox{
					Min: rl.Vector3{X: tx - 1, Y: 0, Z: tz - 1},
					Max: rl.Vector3{X: tx + 1, Y: 10, Z: tz + 1},
				}
				collisionBoxes = append(collisionBoxes, box)
			}
		}
	case Snow:
		for k := 0; k < 2; k++ {
			ix := posX + chunkRand.Float32()*CHUNK_SIZE
			iz := posZ + chunkRand.Float32()*CHUNK_SIZE
			if !isPositionOnRoad(ix-posX, iz-posZ) {
				iglooMesh := rl.GenMeshSphere(5, 16, 16)
				iglooModel := rl.LoadModelFromMesh(iglooMesh)
				iglooTransform := rl.MatrixTranslate(ix, 5, iz)
				iglooModel.Transform = iglooTransform
				setAlbedoColor(&iglooModel, rl.White)
				chunk.Models = append(chunk.Models, iglooModel)
				box := rl.BoundingBox{
					Min: rl.Vector3{X: ix - 5, Y: 0, Z: iz - 5},
					Max: rl.Vector3{X: ix + 5, Y: 10, Z: iz + 5},
				}
				collisionBoxes = append(collisionBoxes, box)
			}
		}
	}

	chunks[coord] = chunk
}

// updateWorld generates new chunks as the player moves.
func updateWorld() {
	playerChunk := getChunkCoord(car.position)
	if playerChunk != lastPlayerChunk {
		dx := playerChunk.X - lastPlayerChunk.X
		dy := playerChunk.Y - lastPlayerChunk.Y
		if dx > 0 {
			for j := playerChunk.Y - 2; j <= playerChunk.Y+2; j++ {
				generateChunk(playerChunk.X+2, j)
			}
		} else if dx < 0 {
			for j := playerChunk.Y - 2; j <= playerChunk.Y+2; j++ {
				generateChunk(playerChunk.X-2, j)
			}
		}
		if dy > 0 {
			for i := playerChunk.X - 2; i <= playerChunk.X+2; i++ {
				generateChunk(i, playerChunk.Y+2)
			}
		} else if dy < 0 {
			for i := playerChunk.X - 2; i <= playerChunk.X+2; i++ {
				generateChunk(i, playerChunk.Y-2)
			}
		}
		lastPlayerChunk = playerChunk
	}
}

// drawWorld renders all chunks in the visible 5x5 grid.
func drawWorld() {
	updateWorld()
	playerChunk := getChunkCoord(car.position)
	for i := playerChunk.X - 2; i <= playerChunk.X+2; i++ {
		for j := playerChunk.Y - 2; j <= playerChunk.Y+2; j++ {
			coord := Coord{i, j}
			if _, exists := chunks[coord]; !exists {
				generateChunk(i, j)
			}
			chunk := chunks[coord]
			for _, model := range chunk.Models {
				rl.DrawModel(model, rl.Vector3{}, 1, rl.White)
			}
		}
	}
}

// initWorld initializes the world by generating a 5x5 grid and resetting collision boxes.
func initWorld() {
	chunks = make(map[Coord]*Chunk)
	collisionBoxes = []rl.BoundingBox{}
	lastPlayerChunk = Coord{0, 0}
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			generateChunk(i, j)
		}
	}
}
