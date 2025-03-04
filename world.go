package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Constants for chunk size and types
const (
	CHUNK_SIZE = 50.0 // Size of each chunk in units
	Highway    = 0
	City       = 1
	Commercial = 2
	Desert     = 3
	Forest     = 4
	Snow       = 5
)

// Coord represents a chunk's position in the world
type Coord struct {
	X, Y int
}

// Chunk holds the type and models for rendering
type Chunk struct {
	Type   int
	Coord  Coord
	Models []rl.Model
}

// Global variables
var (
	chunks map[Coord]*Chunk
	// Adjacency rules: each type lists allowed neighboring types
	allowedNeighbors = [][]int{
		{Highway, City, Commercial, Desert, Forest, Snow}, // Highway
		{Highway, Commercial},                             // City
		{Highway, City},                                   // Commercial
		{Highway, Desert},                                 // Desert
		{Highway, Forest},                                 // Forest
		{Highway, Snow},                                   // Snow
	}
	// Colors for each chunk type's ground
	typeColors = []rl.Color{
		rl.Gray,     // Highway
		rl.DarkGray, // City
		rl.Gray,     // Commercial
		rl.Yellow,   // Desert
		rl.Green,    // Forest
		rl.White,    // Snow
	}

	// Track last player chunk for dynamic generation
	lastPlayerChunk Coord
)

// setAlbedoColor copies the model's material, updates its albedo color,
// and then reassigns the modified material back to the model.
func setAlbedoColor(model *rl.Model, color rl.Color) {
	mat := *model.Materials // copy the material
	mat.Maps.Color = color
	model.Materials = &mat
}

// Initialize the world with an initial 5x5 grid around (0,0)
func initWorld() {
	chunks = make(map[Coord]*Chunk)
	// Generate the initial 5x5 grid around (0,0)
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			generateChunk(i, j)
		}
	}
	// Set the initial player chunk to (0,0)
	lastPlayerChunk = Coord{0, 0}
}

// Convert world position to chunk coordinates
func getChunkCoord(pos rl.Vector3) Coord {
	i := int(math.Floor(float64(pos.X) / CHUNK_SIZE))
	j := int(math.Floor(float64(pos.Z) / CHUNK_SIZE))
	return Coord{i, j}
}

// Determine chunk type based on neighbors
func determineChunkType(i, j int) int {
	neighbors := []Coord{{i - 1, j}, {i + 1, j}, {i, j - 1}, {i, j + 1}}
	existingNeighbors := make([]int, 0)

	// Collect types of existing neighboring chunks
	for _, n := range neighbors {
		if chunk, exists := chunks[n]; exists {
			existingNeighbors = append(existingNeighbors, chunk.Type)
		}
	}

	// If no neighbors, choose randomly
	if len(existingNeighbors) == 0 {
		return rand.Intn(6) // 0 to 5
	}

	// Start with allowed types from the first neighbor
	allowed := make(map[int]bool)
	for _, t := range allowedNeighbors[existingNeighbors[0]] {
		allowed[t] = true
	}

	// Intersect with allowed types from other neighbors
	for _, neighbor := range existingNeighbors[1:] {
		neighborAllowed := allowedNeighbors[neighbor]
		newAllowed := make(map[int]bool)
		for _, t := range neighborAllowed {
			if allowed[t] {
				newAllowed[t] = true
			}
		}
		allowed = newAllowed
	}

	// Convert allowed types to a slice for random selection
	types := make([]int, 0, len(allowed))
	for t := range allowed {
		types = append(types, t)
	}

	// Fallback to Highway if no valid types (shouldn't happen due to Highway's flexibility)
	if len(types) == 0 {
		return Highway
	}

	return types[rand.Intn(len(types))]
}

// Generate a chunk's models based on its type
func generateChunk(i, j int) {
	coord := Coord{i, j}
	// Avoid regenerating an already existing chunk
	if _, exists := chunks[coord]; exists {
		return
	}
	chunk := &Chunk{Type: determineChunkType(i, j), Coord: coord, Models: []rl.Model{}}
	posX := float32(i) * CHUNK_SIZE
	posZ := float32(j) * CHUNK_SIZE

	// Ground plane
	groundMesh := rl.GenMeshPlane(CHUNK_SIZE, CHUNK_SIZE, 1, 1)
	groundModel := rl.LoadModelFromMesh(groundMesh)
	groundTransform := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0, posZ+CHUNK_SIZE/2)
	groundModel.Transform = groundTransform
	setAlbedoColor(&groundModel, typeColors[chunk.Type])
	chunk.Models = append(chunk.Models, groundModel)

	// Define road width
	roadWidth := float32(5.0)

	// "+" shaped road (horizontal and vertical)
	roadMeshH := rl.GenMeshPlane(CHUNK_SIZE, roadWidth, 1, 1)
	roadModelH := rl.LoadModelFromMesh(roadMeshH)
	roadTransformH := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+CHUNK_SIZE/2)
	roadModelH.Transform = roadTransformH
	setAlbedoColor(&roadModelH, rl.DarkGray)
	chunk.Models = append(chunk.Models, roadModelH)

	roadMeshV := rl.GenMeshPlane(roadWidth, CHUNK_SIZE, 1, 1)
	roadModelV := rl.LoadModelFromMesh(roadMeshV)
	roadTransformV := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+CHUNK_SIZE/2)
	roadModelV.Transform = roadTransformV
	setAlbedoColor(&roadModelV, rl.DarkGray)
	chunk.Models = append(chunk.Models, roadModelV)

	// For the central chunk (0,0), only ground and the main road are needed.
	if i == 0 && j == 0 {
		chunks[coord] = chunk
		return
	}

	// Backroads for City and Commercial (forming a square)
	if chunk.Type == City || chunk.Type == Commercial {
		// Top horizontal backroad
		roadModelTop := rl.LoadModelFromMesh(roadMeshH)
		roadTopTransform := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+CHUNK_SIZE-roadWidth)
		roadModelTop.Transform = roadTopTransform
		setAlbedoColor(&roadModelTop, rl.DarkGray)
		chunk.Models = append(chunk.Models, roadModelTop)

		// Bottom horizontal backroad
		roadModelBottom := rl.LoadModelFromMesh(roadMeshH)
		roadBottomTransform := rl.MatrixTranslate(posX+CHUNK_SIZE/2, 0.01, posZ+roadWidth)
		roadModelBottom.Transform = roadBottomTransform
		setAlbedoColor(&roadModelBottom, rl.DarkGray)
		chunk.Models = append(chunk.Models, roadModelBottom)

		// Left vertical backroad
		roadModelLeft := rl.LoadModelFromMesh(roadMeshV)
		roadLeftTransform := rl.MatrixTranslate(posX+roadWidth, 0.01, posZ+CHUNK_SIZE/2)
		roadModelLeft.Transform = roadLeftTransform
		setAlbedoColor(&roadModelLeft, rl.DarkGray)
		chunk.Models = append(chunk.Models, roadModelLeft)

		// Right vertical backroad
		roadModelRight := rl.LoadModelFromMesh(roadMeshV)
		roadRightTransform := rl.MatrixTranslate(posX+CHUNK_SIZE-roadWidth, 0.01, posZ+CHUNK_SIZE/2)
		roadModelRight.Transform = roadRightTransform
		setAlbedoColor(&roadModelRight, rl.DarkGray)
		chunk.Models = append(chunk.Models, roadModelRight)
	}

	// Type-specific objects with seeded randomness
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%d,%d", i, j)))
	seed := h.Sum64()
	chunkRand := rand.New(rand.NewSource(int64(seed)))

	switch chunk.Type {
	case City:
		for k := 0; k < 5; k++ {
			bx := posX + chunkRand.Float32()*CHUNK_SIZE
			bz := posZ + chunkRand.Float32()*CHUNK_SIZE
			if float32(math.Abs(float64(bx-(posX+CHUNK_SIZE/2)))) > roadWidth &&
				float32(math.Abs(float64(bz-(posZ+CHUNK_SIZE/2)))) > roadWidth {
				buildingMesh := rl.GenMeshCube(10, 50, 10)
				buildingModel := rl.LoadModelFromMesh(buildingMesh)
				buildingTransform := rl.MatrixTranslate(bx, 25, bz)
				buildingModel.Transform = buildingTransform
				setAlbedoColor(&buildingModel, rl.Blue)
				chunk.Models = append(chunk.Models, buildingModel)
			}
		}
	case Commercial:
		for k := 0; k < 3; k++ {
			sx := posX + chunkRand.Float32()*CHUNK_SIZE
			sz := posZ + chunkRand.Float32()*CHUNK_SIZE
			if float32(math.Abs(float64(sx-(posX+CHUNK_SIZE/2)))) > roadWidth &&
				float32(math.Abs(float64(sz-(posZ+CHUNK_SIZE/2)))) > roadWidth {
				storeMesh := rl.GenMeshCube(15, 10, 15)
				storeModel := rl.LoadModelFromMesh(storeMesh)
				storeTransform := rl.MatrixTranslate(sx, 5, sz)
				storeModel.Transform = storeTransform
				setAlbedoColor(&storeModel, rl.Purple)
				chunk.Models = append(chunk.Models, storeModel)
			}
		}
	case Desert:
		for k := 0; k < 10; k++ {
			cx := posX + chunkRand.Float32()*CHUNK_SIZE
			cz := posZ + chunkRand.Float32()*CHUNK_SIZE
			cactusMesh := rl.GenMeshCube(1, 5, 1)
			cactusModel := rl.LoadModelFromMesh(cactusMesh)
			cactusTransform := rl.MatrixTranslate(cx, 2.5, cz)
			cactusModel.Transform = cactusTransform
			setAlbedoColor(&cactusModel, rl.Green)
			chunk.Models = append(chunk.Models, cactusModel)
		}
	case Forest:
		for k := 0; k < 20; k++ {
			tx := posX + chunkRand.Float32()*CHUNK_SIZE
			tz := posZ + chunkRand.Float32()*CHUNK_SIZE
			treeMesh := rl.GenMeshCube(2, 10, 2)
			treeModel := rl.LoadModelFromMesh(treeMesh)
			treeTransform := rl.MatrixTranslate(tx, 5, tz)
			treeModel.Transform = treeTransform
			setAlbedoColor(&treeModel, rl.DarkGreen)
			chunk.Models = append(chunk.Models, treeModel)
		}
	case Snow:
		for k := 0; k < 2; k++ {
			ix := posX + chunkRand.Float32()*CHUNK_SIZE
			iz := posZ + chunkRand.Float32()*CHUNK_SIZE
			iglooMesh := rl.GenMeshSphere(5, 16, 16)
			iglooModel := rl.LoadModelFromMesh(iglooMesh)
			iglooTransform := rl.MatrixTranslate(ix, 5, iz)
			iglooModel.Transform = iglooTransform
			setAlbedoColor(&iglooModel, rl.White)
			chunk.Models = append(chunk.Models, iglooModel)
		}
	}

	chunks[coord] = chunk
}

// Update the world by generating new chunks if the player has moved into a new chunk.
func updateWorld() {
	playerChunk := getChunkCoord(car.position)
	if playerChunk != lastPlayerChunk {
		dx := playerChunk.X - lastPlayerChunk.X
		dy := playerChunk.Y - lastPlayerChunk.Y

		// If moved horizontally, generate the new column on the side the player moved to.
		if dx > 0 {
			// Moved east; generate column at player's east edge (playerChunk.X + 2)
			for j := playerChunk.Y - 2; j <= playerChunk.Y+2; j++ {
				generateChunk(playerChunk.X+2, j)
			}
		} else if dx < 0 {
			// Moved west; generate column at player's west edge (playerChunk.X - 2)
			for j := playerChunk.Y - 2; j <= playerChunk.Y+2; j++ {
				generateChunk(playerChunk.X-2, j)
			}
		}
		// If moved vertically, generate the new row on the side the player moved to.
		if dy > 0 {
			// Moved south; generate row at player's south edge (playerChunk.Y + 2)
			for i := playerChunk.X - 2; i <= playerChunk.X+2; i++ {
				generateChunk(i, playerChunk.Y+2)
			}
		} else if dy < 0 {
			// Moved north; generate row at player's north edge (playerChunk.Y - 2)
			for i := playerChunk.X - 2; i <= playerChunk.X+2; i++ {
				generateChunk(i, playerChunk.Y-2)
			}
		}
		lastPlayerChunk = playerChunk
	}
}

// Draw the visible 5x5 chunk grid
func drawWorld() {
	// Update the world before drawing
	updateWorld()

	playerChunk := getChunkCoord(car.position)
	for i := playerChunk.X - 2; i <= playerChunk.X+2; i++ {
		for j := playerChunk.Y - 2; j <= playerChunk.Y+2; j++ {
			coord := Coord{i, j}
			// If chunk does not exist (should be generated by updateWorld) generate it
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
