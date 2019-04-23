package main

import "math/rand"

type bufferedUniverse struct {
	generation    uint
	rows, columns int
	cells         [2][]bool
}

type Universe interface {
	Generation() uint
	Iterate()
	RowCount() int
	ColumnCount() int
	IsAlive(row, column int) bool
}

func (b *bufferedUniverse) Generation() uint {
	return b.generation
}

func (b *bufferedUniverse) Iterate() {
	b.generation++
}

func (b *bufferedUniverse) RowCount() int {
	return b.rows
}

func (b *bufferedUniverse) ColumnCount() int {
	return b.columns
}

func (b *bufferedUniverse) IsAlive(row, column int) bool {
	return false
}

func NewBufferedUniverse(rows, columns int, random *rand.Rand) *bufferedUniverse {
	size := rows * columns
	cellBuffer := [2][]bool{
		make([]bool, size, size),
		make([]bool, size, size),
	}

	for i := 0; i < size; i++ {
		cellBuffer[0][i] = random.Intn(2) == 0
	}

	return &bufferedUniverse{
		generation: uint(0),
		rows:       rows,
		columns:    columns,
		cells:      cellBuffer,
	}
}
