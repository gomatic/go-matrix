package matrix

import "iter"

// Width is the number of columns in a [Matrix].
type Width int

// Height is the number of rows in a [Matrix].
type Height int

// Point is a cell coordinate. The origin is the top-left, X increases to the
// right, and Y increases downward.
type Point struct {
	X int
	Y int
}

// Fill produces the initial value for the cell at a coordinate. A nil Fill
// passed to [New] leaves every cell at its zero value.
type Fill[T any] func(Point) T

// Transform maps a cell's coordinate and current value to its new value. It is
// applied by [Matrix.Map] to produce a new grid.
type Transform[T any] func(Point, T) T

// Matrix is an immutable rectangular grid of cells of type T. The zero Matrix is
// a valid, empty grid that contains no cells.
type Matrix[T any] struct {
	data   []T
	width  int
	height int
}

// New builds a width×height grid, setting each cell from fill. It returns
// [ErrNonPositiveWidth] or [ErrNonPositiveHeight] when a dimension is not
// greater than zero. A nil fill leaves every cell at its zero value.
func New[T any](width Width, height Height, fill Fill[T]) (Matrix[T], error) {
	if width <= 0 {
		return Matrix[T]{}, ErrNonPositiveWidth
	}
	if height <= 0 {
		return Matrix[T]{}, ErrNonPositiveHeight
	}
	w, h := int(width), int(height)
	data := make([]T, w*h)
	if fill != nil {
		for i := range data {
			data[i] = fill(Point{X: i % w, Y: i / w})
		}
	}
	return Matrix[T]{data: data, width: w, height: h}, nil
}

// Width returns the number of columns.
func (m Matrix[T]) Width() Width { return Width(m.width) }

// Height returns the number of rows.
func (m Matrix[T]) Height() Height { return Height(m.height) }

// Contains reports whether p lies within the grid.
func (m Matrix[T]) Contains(p Point) bool {
	return p.X >= 0 && p.X < m.width && p.Y >= 0 && p.Y < m.height
}

// At returns the value at p and whether p is within the grid. When p is out of
// bounds it returns the zero value and false.
func (m Matrix[T]) At(p Point) (T, bool) {
	if !m.Contains(p) {
		var zero T
		return zero, false
	}
	return m.data[m.index(p)], true
}

// Cells iterates every cell in row-major order, yielding each coordinate and its
// value.
func (m Matrix[T]) Cells() iter.Seq2[Point, T] {
	return func(yield func(Point, T) bool) {
		for i, value := range m.data {
			if !yield(m.point(i), value) {
				return
			}
		}
	}
}

// Map returns a new grid of the same size whose cells are the results of
// applying transform to each cell of m. The receiver is unchanged.
func (m Matrix[T]) Map(transform Transform[T]) Matrix[T] {
	data := make([]T, len(m.data))
	for i, value := range m.data {
		data[i] = transform(m.point(i), value)
	}
	return Matrix[T]{data: data, width: m.width, height: m.height}
}

// index returns the row-major offset of an in-bounds point.
func (m Matrix[T]) index(p Point) int { return p.Y*m.width + p.X }

// point returns the coordinate of a row-major offset.
func (m Matrix[T]) point(index int) Point {
	return Point{X: index % m.width, Y: index / m.width}
}
