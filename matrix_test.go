package matrix

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
)

// collect drains an iterator into parallel coordinate and value slices.
func collect[T any](seq iter.Seq2[Point, T]) ([]Point, []T) {
	var points []Point
	var values []T
	for p, v := range seq {
		points = append(points, p)
		values = append(values, v)
	}
	return points, values
}

// index encodes a coordinate as a distinct cell value for assertions.
func index(width int) Fill[int] {
	return func(p Point) int { return p.Y*width + p.X }
}

func TestNew(t *testing.T) {
	tests := []struct {
		wantErr error
		name    string
		width   Width
		height  Height
	}{
		{name: "valid", width: 3, height: 2},
		{name: "zero width", width: 0, height: 2, wantErr: ErrNonPositiveWidth},
		{name: "negative width", width: -1, height: 2, wantErr: ErrNonPositiveWidth},
		{name: "zero height", width: 3, height: 0, wantErr: ErrNonPositiveHeight},
		{name: "negative height", width: 3, height: -4, wantErr: ErrNonPositiveHeight},
		{name: "both non-positive reports the width first", width: 0, height: 0, wantErr: ErrNonPositiveWidth},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			try := assert.New(t)
			m, err := New(tt.width, tt.height, index(int(tt.width)))
			if tt.wantErr != nil {
				try.ErrorIs(err, tt.wantErr)
				return
			}
			try.NoError(err)
			try.Equal(tt.width, m.Width())
			try.Equal(tt.height, m.Height())
		})
	}
}

func TestNew_fill(t *testing.T) {
	try := assert.New(t)
	m, err := New(3, 2, func(p Point) Point { return p })
	try.NoError(err)
	v, ok := m.At(Point{X: 2, Y: 1})
	try.True(ok)
	try.Equal(Point{X: 2, Y: 1}, v)
}

func TestNew_nilFill(t *testing.T) {
	try := assert.New(t)
	m, err := New[int](2, 2, nil)
	try.NoError(err)
	v, ok := m.At(Point{X: 1, Y: 1})
	try.True(ok)
	try.Equal(0, v)
}

func TestContainsAndAt(t *testing.T) {
	try := assert.New(t)
	m, _ := New(3, 2, index(3))

	try.True(m.Contains(Point{X: 2, Y: 1}))
	v, ok := m.At(Point{X: 2, Y: 1})
	try.True(ok)
	try.Equal(5, v)

	for _, p := range []Point{{X: -1, Y: 0}, {X: 3, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 2}} {
		try.False(m.Contains(p))
		v, ok := m.At(p)
		try.False(ok)
		try.Equal(0, v)
	}
}

func TestCells(t *testing.T) {
	try := assert.New(t)
	m, _ := New(2, 2, index(2))
	points, values := collect(m.Cells())
	try.Equal([]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}, points)
	try.Equal([]int{0, 1, 2, 3}, values)
}

func TestCells_break(t *testing.T) {
	try := assert.New(t)
	m, _ := New(2, 2, index(2))
	count := 0
	for range m.Cells() {
		count++
		break
	}
	try.Equal(1, count)
}

func TestMap(t *testing.T) {
	try := assert.New(t)
	m, _ := New(2, 2, index(2))

	scaled := m.Map(func(_ Point, v int) int { return v * 10 })
	v, ok := scaled.At(Point{X: 1, Y: 1})
	try.True(ok)
	try.Equal(30, v)

	// The receiver is unchanged and the new grid keeps the dimensions.
	original, _ := m.At(Point{X: 1, Y: 1})
	try.Equal(3, original)
	try.Equal(m.Width(), scaled.Width())
	try.Equal(m.Height(), scaled.Height())

	// Transform receives each cell's coordinate.
	coords := m.Map(func(p Point, _ int) int { return p.X })
	x, _ := coords.At(Point{X: 1, Y: 0})
	try.Equal(1, x)
}
