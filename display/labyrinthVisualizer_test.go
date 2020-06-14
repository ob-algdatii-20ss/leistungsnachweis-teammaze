package display

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ob-algdatii-20ss/leistungsnachweis-teammaze/common"
)

func testingCubeConstructor(x, y, z, xSize, ySize, zSize float32) Cube {
	return newCube(x, y, z, xSize, ySize, zSize, getTestRenderInfo())
}

func getTestRenderInfo() *renderInfo {
	return &renderInfo{}
}

func TestMakeConnectionFailsOnNonAdjacent(t *testing.T) {
	loc1 := common.NewLocation(0, 0, 0)
	loc2 := common.NewLocation(1, 1, 1)

	defer func() {
		want := "tried to make connection between non-adjacent locations"
		if got := recover(); got != nil {
			if got != want {
				t.Errorf("expected \"%s\" but got \"%s\"", want, got)
			}
		} else {
			t.Error("Expected panic, got none")
		}
	}()

	makeConnection(loc1, loc2, testingCubeConstructor)
}

func TestMakeConnectionHasCorrectCenter(t *testing.T) {
	loc1 := common.NewLocation(1, 2, 1)
	loc2 := common.NewLocation(1, 1, 1)
	cube := makeConnection(loc1, loc2, testingCubeConstructor)

	got := cube.Transform.translation.Mul4x1(mgl32.Vec4{0, 0, 0, 1}).Vec3()
	want := mgl32.Vec3{1, 1.5, 1}

	if got != want {
		t.Errorf("got: %v\nexpected: %v", got, want)
	}
}

func TestMakeConnectionHasCorrectScale(t *testing.T) {
	loc1 := common.NewLocation(3, 3, 3)
	loc2 := common.NewLocation(4, 3, 3)
	cube := makeConnection(loc1, loc2, testingCubeConstructor)

	got := cube.Transform.scale.Mul4x1(mgl32.Vec4{1, 1, 1, 1}).Vec3()
	want := mgl32.Vec3{0.5, 0.25, 0.25}

	if got != want {
		t.Errorf("got: %v\nexpected: %v", got, want)
	}
}

func TestCheckAndMake(t *testing.T) {
	maxLoc := common.NewLocation(2, 2, 2)
	lab := common.NewLabyrinth(maxLoc)

	baseLoc := common.NewLocation(1, 1, 1)
	lab.Connect(baseLoc, common.NewLocation(1, 2, 1))
	lab.Connect(baseLoc, common.NewLocation(1, 1, 2))

	wantedCubes := []Cube{
		newCube(1, 1.5, 1, 0.25, 0.5, 0.25, getTestRenderInfo()),
		newCube(1, 1, 1.5, 0.25, 0.25, 0.5, getTestRenderInfo()),
	}
	cubes := makeConnections(&lab, baseLoc, testingCubeConstructor)

	compareCubeSlices(t, cubes, wantedCubes)
}

func TestExploreLabyrinth(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	wantedCubes := []Cube{
		// Nodes
		newCube(0, 0, 0, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(0, 0, 1, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(0, 0, 2, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(0, 1, 0, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(0, 1, 1, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(0, 1, 2, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(1, 0, 0, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(1, 0, 1, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(1, 0, 2, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(1, 1, 0, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(1, 1, 1, 0.5, 0.5, 0.5, getTestRenderInfo()),
		newCube(1, 1, 2, 0.5, 0.5, 0.5, getTestRenderInfo()),
		// Connections
		newCube(0, 0.5, 0, 0.25, 0.5, 0.25, getTestRenderInfo()),
		newCube(0, 0, 0.5, 0.25, 0.25, 0.5, getTestRenderInfo()),
		newCube(1, 0, 1.5, 0.25, 0.25, 0.5, getTestRenderInfo()),
		newCube(1, 0.5, 1, 0.25, 0.5, 0.25, getTestRenderInfo()),
		newCube(0, 0.5, 1, 0.25, 0.5, 0.25, getTestRenderInfo()),
		newCube(1, 1, 1.5, 0.25, 0.25, 0.5, getTestRenderInfo()),
	}

	lab.Connect(common.NewLocation(0, 0, 0), common.NewLocation(0, 1, 0))
	lab.Connect(common.NewLocation(0, 0, 0), common.NewLocation(0, 0, 1))
	lab.Connect(common.NewLocation(1, 0, 1), common.NewLocation(1, 0, 2))
	lab.Connect(common.NewLocation(1, 0, 1), common.NewLocation(1, 1, 1))
	lab.Connect(common.NewLocation(0, 0, 1), common.NewLocation(0, 1, 1))
	lab.Connect(common.NewLocation(1, 1, 1), common.NewLocation(1, 1, 2))

	cubes := exploreLabyrinth(&lab, testingCubeConstructor)

	compareCubeSlices(t, cubes, wantedCubes)
}

func TestNewLabyrinthVisualizerPanicsOnNil(t *testing.T) {
	defer func() {
		want := "passed labyrinth has to be valid"
		if got := recover(); got != nil {
			if got != want {
				t.Errorf("Unexpected panic: \"%s\", expected: \"%s\"", got, want)
			}
		} else {
			t.Errorf("Expected panic, got none")
		}
	}()

	NewLabyrinthVisualizer(nil, nil, nil, testingCubeConstructor)
}

func TestStepsMayBeNilInConstructor(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	NewLabyrinthVisualizer(&lab, nil, nil, testingCubeConstructor)
}

func TestDoStepPanicsWhenStepNil(t *testing.T) {
	defer func() {
		want := "cannot do step: steps is nil"
		if got := recover(); got != nil {
			if got != want {
				t.Errorf("Unexpected panic: \"%s\", expected: \"%s\"", got, want)
			}
		} else {
			t.Errorf("Expected panic, got none")
		}
	}()

	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	vis := NewLabyrinthVisualizer(&lab, nil, nil, testingCubeConstructor)
	vis.DoStep()
}

func TestDoStepPanicsWhenColorConverterNil(t *testing.T) {
	defer func() {
		want := "cannot do step: color converter is nil"
		if got := recover(); got != nil {
			if got != want {
				t.Errorf("Unexpected panic: \"%s\", expected: \"%s\"", got, want)
			}
		} else {
			t.Errorf("Expected panic, got none")
		}
	}()

	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	vis := NewLabyrinthVisualizer(&lab, make([]common.Pair, 0), nil, testingCubeConstructor)
	vis.DoStep()
}

type testConverter struct {
	mapping map[string]mgl32.Vec4
}

func newTestConverter(mapping map[string]mgl32.Vec4) testConverter {
	return testConverter{
		mapping: mapping,
	}
}

func (t testConverter) StepToColor(step common.Pair, cubes []Cube) (*Cube, mgl32.Vec4) {
	location := step.GetFirst().(common.Location)
	symbol := step.GetSecond().(string)

	x, y, z := location.As3DCoordinates()
	location3D := mgl32.Vec3{float32(x), float32(y), float32(z)}

	for index, cube := range cubes {
		if cube.Transform.GetTranslation() == location3D {
			return &cubes[index], t.mapping[symbol]
		}
	}

	return nil, mgl32.Vec4{}
}

func (t testConverter) ColorMap() map[string]mgl32.Vec4 {
	return t.mapping
}

func TestDoStepColorsFirstCubeOnCall(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	steps := []common.Pair{common.NewPair(lab.GetMaxLocation(), "START")}

	mapping := map[string]mgl32.Vec4{
		"START": {1, 1, 1, 1},
	}
	vis := NewLabyrinthVisualizer(&lab, steps, newTestConverter(mapping), testingCubeConstructor)
	vis.DoStep()

	for _, cube := range vis.cubes {
		if cube.Transform.GetTranslation() == (mgl32.Vec3{1, 1, 2}) {
			if cube.info.color != (mgl32.Vec4{1, 1, 1, 1}) {
				t.Errorf("Expecting color of %v to be %v, but was %v", cube, mgl32.Vec4{1, 1, 1, 1}, cube.info.color)
			}
		}
	}
}

func TestDoStepColorsSecondCubeOnTwoCalls(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	steps := []common.Pair{
		common.NewPair(common.NewLocation(0, 0, 0), "START"),
		common.NewPair(common.NewLocation(0, 1, 0), "ADD"),
	}

	mapping := map[string]mgl32.Vec4{
		"START": {1, 1, 1, 1},
		"ADD":   {0, 1, 0, 1},
	}

	vis := NewLabyrinthVisualizer(&lab, steps, newTestConverter(mapping), testingCubeConstructor)
	vis.DoStep()
	vis.DoStep()

	for _, cube := range vis.cubes {
		if cube.Transform.GetTranslation() == (mgl32.Vec3{0, 1, 0}) {
			if cube.info.color != (mgl32.Vec4{0, 1, 0, 1}) {
				t.Errorf("Expecting color of %v to be %v, but was %v", cube, mgl32.Vec4{0, 1, 0, 1}, cube.info.color)
			}
		}
	}
}

func TestDoStepColorsConnectionToSecondCubeOnTwoCalls(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	lab.Connect(common.NewLocation(0, 0, 0), common.NewLocation(0, 1, 0))

	steps := []common.Pair{
		common.NewPair(common.NewLocation(0, 0, 0), "START"),
		common.NewPair(common.NewLocation(0, 1, 0), "ADD"),
	}

	mapping := map[string]mgl32.Vec4{
		"START": {1, 1, 1, 1},
		"ADD":   {0, 1, 0, 1},
	}

	vis := NewLabyrinthVisualizer(&lab, steps, newTestConverter(mapping), testingCubeConstructor)
	vis.DoStep()
	vis.DoStep()

	requireCube(t, vis.cubes, locationCondition(mgl32.Vec3{0, 0.5, 0}), colorCondition(mgl32.Vec4{1, 1, 1, 1}),
		func(cube Cube) string {
			return fmt.Sprintf("Expected Cube at %v to have color %v but had %v",
				cube.Transform.GetTranslation(), mgl32.Vec4{0, 1, 0, 1}, cube.info.color)
		})
}

func TestDoStepColorsConnectionAfterFirstCube(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	lab.Connect(common.NewLocation(0, 0, 0), common.NewLocation(0, 1, 0))
	lab.Connect(common.NewLocation(0, 1, 0), common.NewLocation(0, 1, 1))

	steps := []common.Pair{
		common.NewPair(common.NewLocation(0, 0, 0), "START"),
		common.NewPair(common.NewLocation(0, 1, 0), "ADD"),
		common.NewPair(common.NewLocation(0, 1, 1), "END"),
	}

	mapping := map[string]mgl32.Vec4{
		"START": {1, 1, 1, 1},
		"ADD":   {0, 1, 0, 1},
		"END":   {0, 0, 0, 0},
	}

	vis := NewLabyrinthVisualizer(&lab, steps, newTestConverter(mapping), testingCubeConstructor)
	vis.DoStep()
	vis.DoStep()
	vis.DoStep()

	//Connection 0 -> 1 has color of 0
	requireCube(t, vis.cubes, locationCondition(mgl32.Vec3{0, 0.5, 0}), colorCondition(mgl32.Vec4{1, 1, 1, 1}),
		func(cube Cube) string {
			return fmt.Sprintf("Expected Cube at %v to have color %v but had %v",
				cube.Transform.GetTranslation(), mgl32.Vec4{0, 1, 0, 1}, cube.info.color)
		})

	//Connection 1 -> 2 has color of 1
	requireCube(t, vis.cubes, locationCondition(mgl32.Vec3{0, 1, 0.5}), colorCondition(mgl32.Vec4{0, 1, 0, 1}),
		func(cube Cube) string {
			return fmt.Sprintf("Expected Cube at %v to have color %v but had %v",
				cube.Transform.GetTranslation(), mgl32.Vec4{0, 1, 0, 1}, cube.info.color)
		})
}

func TestDoStepNonAdjacent(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	lab.Connect(common.NewLocation(0, 0, 0), common.NewLocation(0, 1, 0))
	lab.Connect(common.NewLocation(0, 1, 0), common.NewLocation(0, 1, 1))

	steps := []common.Pair{
		common.NewPair(common.NewLocation(0, 0, 0), "START"),
		common.NewPair(common.NewLocation(0, 1, 0), "ADD"),
		common.NewPair(common.NewLocation(0, 1, 2), "END"),
	}

	mapping := map[string]mgl32.Vec4{
		"START": {1, 1, 1, 1},
		"ADD":   {0, 1, 0, 1},
		"END":   {0, 0, 0, 0},
	}

	vis := NewLabyrinthVisualizer(&lab, steps, newTestConverter(mapping), testingCubeConstructor)
	vis.DoStep()
	vis.DoStep()
	vis.DoStep()

	//Connection from 1 has default color
	requireCube(t, vis.cubes, locationCondition(mgl32.Vec3{0, 1, 0.5}), colorCondition(defaultCubeColor()),
		func(cube Cube) string {
			return fmt.Sprintf("Expected Cube at %v to have color %v but had %v",
				cube.Transform.GetTranslation(), defaultCubeColor(), cube.info.color)
		})

	//End Cube has correct color
	requireCube(t, vis.cubes, locationCondition(mgl32.Vec3{0, 1, 2}), colorCondition(mgl32.Vec4{0, 0, 0, 0}),
		func(cube Cube) string {
			return fmt.Sprintf("Expected Cube at %v to have color %v but had %v",
				cube.Transform.GetTranslation(), mgl32.Vec4{0, 0, 0, 0}, cube.info.color)
		})
}

func TestDoStepAdjacentOutOfOrder(t *testing.T) {
	maxLoc := common.NewLocation(1, 1, 2)
	lab := common.NewLabyrinth(maxLoc)

	lab.Connect(common.NewLocation(0, 0, 0), common.NewLocation(0, 1, 0))
	lab.Connect(common.NewLocation(0, 1, 0), common.NewLocation(0, 1, 1))

	steps := []common.Pair{
		common.NewPair(common.NewLocation(0, 0, 0), "START"),
		common.NewPair(common.NewLocation(0, 1, 1), "END"),
		common.NewPair(common.NewLocation(0, 1, 0), "ADD"),
	}

	mapping := map[string]mgl32.Vec4{
		"START": {1, 1, 1, 1},
		"ADD":   {0, 1, 0, 1},
		"END":   {0, 0, 0, 0},
	}

	vis := NewLabyrinthVisualizer(&lab, steps, newTestConverter(mapping), testingCubeConstructor)
	vis.DoStep()
	vis.DoStep()
	vis.DoStep()

	//Connection 0 -> 1 has color of 0
	requireCube(t, vis.cubes, locationCondition(mgl32.Vec3{0, 0.5, 0}), colorCondition(mgl32.Vec4{1, 1, 1, 1}),
		func(cube Cube) string {
			return fmt.Sprintf("Expected Cube at %v to have color %v but had %v",
				cube.Transform.GetTranslation(), mgl32.Vec4{1, 1, 1, 1}, cube.info.color)
		})

	//Connection 1 -> 2 has color of 1
	requireCube(t, vis.cubes, locationCondition(mgl32.Vec3{0, 1, 0.5}), colorCondition(mgl32.Vec4{0, 1, 0, 1}),
		func(cube Cube) string {
			return fmt.Sprintf("Expected Cube at %v to have color %v but had %v",
				cube.Transform.GetTranslation(), mgl32.Vec4{0, 1, 0, 1}, cube.info.color)
		})
}

type cubeCondition func(Cube) bool

func locationCondition(location mgl32.Vec3) cubeCondition {
	return func(cube Cube) bool {
		return cube.Transform.GetTranslation() == location
	}
}

func colorCondition(color mgl32.Vec4) cubeCondition {
	return func(cube Cube) bool {
		return cube.info.color == color
	}
}

func requireCube(t *testing.T, cubes []Cube, isCube cubeCondition, condition cubeCondition, error func(Cube) string) {
	t.Helper()

	cubeFound := false

	for _, cube := range cubes {
		if isCube(cube) {
			if !condition(cube) {
				t.Errorf(error(cube))
			}

			cubeFound = true
		}
	}

	if !cubeFound {
		t.Errorf("Cube not found")
	}
}

// Helpers:

func compareCubeSlices(t *testing.T, cubes []Cube, wantedCubes []Cube) {
	t.Helper()

	if len(cubes) < len(wantedCubes) {
		t.Errorf("Not enough cubes")
	} else if len(cubes) > len(wantedCubes) {
		t.Errorf("Too many cubes")
	}

	for _, cube := range cubes {
		isInWanted := false

		for i, wantedCube := range wantedCubes {
			if cube.Transform == wantedCube.Transform {
				wantedCubes = append(wantedCubes[0:i], wantedCubes[i+1:]...)
				isInWanted = true

				break
			}
		}

		if !isInWanted {
			t.Errorf("Unexpected Cube at %v", cube.Transform.GetTranslation())
		}
	}

	if len(wantedCubes) != 0 {
		for _, wantedCube := range wantedCubes {
			t.Errorf("Failed to make cube at %v\n", wantedCube.Transform.GetTranslation())
		}
	}
}
