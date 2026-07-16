package matrix

import "iter"

// Distance is a Chebyshev (chessboard) distance measured in cells, used by
// [Matrix.Ring].
type Distance int

// Radius is a circle radius measured in cells, used by [Matrix.Circle].
type Radius int

// Neighbors iterates the up-to-eight cells surrounding center, clockwise from
// the top-left, skipping any that fall outside the grid.
func (m Matrix[T]) Neighbors(center Point) iter.Seq2[Point, T] {
	return m.seq(ringPoints(center, 1))
}

// Ring iterates the cells on the square ring at the given Chebyshev distance
// from center, clockwise from the top-left corner, skipping any outside the
// grid. A distance of zero or less yields nothing.
func (m Matrix[T]) Ring(center Point, distance Distance) iter.Seq2[Point, T] {
	return m.seq(ringPoints(center, distance))
}

// Circle iterates the cells on the Bresenham circle of the given radius around
// center, each cell at most once, skipping any outside the grid. A radius of
// zero or less yields nothing.
func (m Matrix[T]) Circle(center Point, radius Radius) iter.Seq2[Point, T] {
	return m.seq(circlePoints(center, radius))
}

// seq iterates the given points, yielding each in-bounds coordinate and its
// value in order.
func (m Matrix[T]) seq(points []Point) iter.Seq2[Point, T] {
	return func(yield func(Point, T) bool) {
		for _, p := range points {
			value, ok := m.At(p)
			if ok && !yield(p, value) {
				return
			}
		}
	}
}

// ringPoints returns the perimeter of the square ring at the given Chebyshev
// distance from center, clockwise from the top-left corner.
func ringPoints(center Point, distance Distance) []Point {
	if distance <= 0 {
		return nil
	}
	cx, cy, d := center.X, center.Y, int(distance)
	points := make([]Point, 0, 8*d)
	for x := cx - d; x < cx+d; x++ {
		points = append(points, Point{X: x, Y: cy - d})
	}
	for y := cy - d; y < cy+d; y++ {
		points = append(points, Point{X: cx + d, Y: y})
	}
	for x := cx + d; x > cx-d; x-- {
		points = append(points, Point{X: x, Y: cy + d})
	}
	for y := cy + d; y > cy-d; y-- {
		points = append(points, Point{X: cx - d, Y: y})
	}
	return points
}

// circlePoints returns the distinct cells on the Bresenham circle of the given
// radius around center, in emission order.
func circlePoints(center Point, radius Radius) []Point {
	if radius <= 0 {
		return nil
	}
	return dedupe(bresenhamOctants(center, radius))
}

// bresenhamOctants returns the eight octant-symmetric points of each step of the
// midpoint circle algorithm; the result may contain duplicates.
func bresenhamOctants(center Point, radius Radius) []Point {
	cx, cy := center.X, center.Y
	points := make([]Point, 0, 8*radius)
	x, y, err := int(radius), 0, 1-int(radius)
	for x >= y {
		points = append(points,
			Point{X: cx + x, Y: cy + y}, Point{X: cx + y, Y: cy + x},
			Point{X: cx - y, Y: cy + x}, Point{X: cx - x, Y: cy + y},
			Point{X: cx - x, Y: cy - y}, Point{X: cx - y, Y: cy - x},
			Point{X: cx + y, Y: cy - x}, Point{X: cx + x, Y: cy - y})
		y++
		if err < 0 {
			err += 2*y + 1
		} else {
			x--
			err += 2*(y-x) + 1
		}
	}
	return points
}

// dedupe returns points with duplicates removed, preserving first-seen order.
func dedupe(points []Point) []Point {
	seen := make(map[Point]struct{}, len(points))
	unique := make([]Point, 0, len(points))
	for _, p := range points {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		unique = append(unique, p)
	}
	return unique
}
