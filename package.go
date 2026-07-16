// Package matrix is an immutable, generic 2-D grid.
//
// A [Matrix] holds width×height cells of an arbitrary type, addressed by [Point]
// coordinates whose origin is the top-left with X increasing to the right and Y
// downward. A grid is built whole by [New] and never mutated in place —
// [Matrix.Map] returns a new grid — so a value is safe to copy, alias, and read
// from multiple goroutines without locking.
//
// Alongside direct access ([Matrix.At], [Matrix.Cells]) it offers the common
// grid traversals — [Matrix.Neighbors], [Matrix.Ring], and [Matrix.Circle] — as
// range-over-func iterators that visit only the in-bounds cells they cross.
package matrix
