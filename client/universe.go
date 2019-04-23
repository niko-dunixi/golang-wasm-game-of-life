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
	bufferIndex := b.currentBufferIndex()
	row = forceInRange(row, b.RowCount())
	column = forceInRange(column, b.ColumnCount())
	index := asIndex(b.columns, row, column)
	isAlive := b.cells[bufferIndex][index]
	return isAlive
}

func (b *bufferedUniverse) setNextLife(row, column int, life bool) {
	bufferIndex := b.nextBufferIndex()
	row = forceInRange(row, b.RowCount())
	column = forceInRange(column, b.ColumnCount())
	index := asIndex(b.columns, row, column)
	b.cells[bufferIndex][index] = life
}

func (b *bufferedUniverse) currentBufferIndex() int {
	return int(b.generation % 2)
}

func (b *bufferedUniverse) nextBufferIndex() int {
	return int((b.generation + 1) % 2)
}

func asIndex(columnCount, row, column int) int {
	return columnCount*row + column
}

func forceInRange(value, maxValue int) int {
	positiveValue := uint(value)
	positiveMaxValue := uint(maxValue)
	return int(positiveValue % positiveMaxValue)
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
