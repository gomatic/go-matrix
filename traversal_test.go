package matrix

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNeighbors_center(t *testing.T) {
	try := assert.New(t)
	m, _ := New(5, 5, index(5))
	points, _ := collect(m.Neighbors(Point{X: 2, Y: 2}))
	try.Equal([]Point{
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 3, Y: 1},
		{X: 3, Y: 2},
		{X: 3, Y: 3},
		{X: 2, Y: 3},
		{X: 1, Y: 3},
		{X: 1, Y: 2},
	}, points)
}

func TestNeighbors_corner(t *testing.T) {
	try := assert.New(t)
	m, _ := New(5, 5, index(5))
	// The out-of-grid neighbours are skipped, leaving the three in-bounds cells.
	points, _ := collect(m.Neighbors(Point{X: 0, Y: 0}))
	try.Equal([]Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}}, points)
}

func TestRing_nonPositive(t *testing.T) {
	try := assert.New(t)
	m, _ := New(5, 5, index(5))
	points, _ := collect(m.Ring(Point{X: 2, Y: 2}, 0))
	try.Empty(points)
	points, _ = collect(m.Ring(Point{X: 2, Y: 2}, -1))
	try.Empty(points)
}

func TestRing_matchesNeighborsAtOne(t *testing.T) {
	try := assert.New(t)
	m, _ := New(5, 5, index(5))
	ring, _ := collect(m.Ring(Point{X: 2, Y: 2}, 1))
	neighbors, _ := collect(m.Neighbors(Point{X: 2, Y: 2}))
	try.Equal(neighbors, ring)
}

func TestRing_distanceTwo(t *testing.T) {
	try := assert.New(t)
	m, _ := New(5, 5, index(5))
	points, values := collect(m.Ring(Point{X: 2, Y: 2}, 2))
	try.Len(points, 16)                     // the full 8*distance perimeter is in bounds
	try.Equal(Point{X: 0, Y: 0}, points[0]) // clockwise from the top-left corner
	try.Equal(0, values[0])
}

func TestCircle_nonPositive(t *testing.T) {
	try := assert.New(t)
	m, _ := New(11, 11, index(11))
	points, _ := collect(m.Circle(Point{X: 5, Y: 5}, 0))
	try.Empty(points)
	points, _ = collect(m.Circle(Point{X: 5, Y: 5}, -3))
	try.Empty(points)
}

func TestCircle_radiusOne(t *testing.T) {
	try := assert.New(t)
	m, _ := New(11, 11, index(11))
	points, _ := collect(m.Circle(Point{X: 5, Y: 5}, 1))
	// The four axis-adjacent cells, each yielded once (duplicates removed).
	try.ElementsMatch([]Point{
		{X: 6, Y: 5}, {X: 5, Y: 6}, {X: 4, Y: 5}, {X: 5, Y: 4},
	}, points)
	try.Len(points, distinct(points))
}

func TestCircle_onOutline(t *testing.T) {
	try := assert.New(t)
	m, _ := New(11, 11, index(11))
	center := Point{X: 5, Y: 5}
	points, _ := collect(m.Circle(center, 4))
	try.NotEmpty(points)
	try.Len(points, distinct(points)) // no duplicates survive
	for _, p := range points {
		dx, dy := float64(p.X-center.X), float64(p.Y-center.Y)
		try.Equal(4.0, math.Round(math.Hypot(dx, dy))) // every cell lies on the radius-4 outline
	}
}

func TestCircle_clipsToGrid(t *testing.T) {
	try := assert.New(t)
	m, _ := New(5, 5, index(5))
	// A circle centred on a corner: only the in-bounds arc is yielded.
	points, _ := collect(m.Circle(Point{X: 0, Y: 0}, 3))
	try.NotEmpty(points)
	for _, p := range points {
		try.True(m.Contains(p))
	}
}

func TestCircle_break(t *testing.T) {
	try := assert.New(t)
	m, _ := New(11, 11, index(11))
	count := 0
	for range m.Circle(Point{X: 5, Y: 5}, 3) {
		count++
		break
	}
	try.Equal(1, count)
}

// distinct returns the number of unique points, for duplicate assertions.
func distinct(points []Point) int {
	seen := make(map[Point]struct{}, len(points))
	for _, p := range points {
		seen[p] = struct{}{}
	}
	return len(seen)
}
