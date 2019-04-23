package main

import (
	"bytes"
	"math/rand"
)

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
	for currentRow := 0; currentRow < b.RowCount(); currentRow++ {
		for currentColumn := 0; currentColumn < b.ColumnCount(); currentColumn++ {
			// https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life#Rules
			//    Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			//    Any live cell with two or three live neighbours lives on to the next generation.
			//    Any live cell with more than three live neighbours dies, as if by overpopulation.
			//    Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			currentlyAlive := b.IsAlive(currentRow, currentColumn)
			currentlyDead := !currentlyAlive
			liveNeighborCount := b.countLiveNeighbors(currentRow, currentColumn)
			if currentlyAlive && liveNeighborCount < 2 {
				b.setNextLife(currentRow, currentColumn, false)
			} else if currentlyAlive && (liveNeighborCount >= 2 || liveNeighborCount <= 3) {
				b.setNextLife(currentRow, currentColumn, true)
			} else if currentlyAlive && liveNeighborCount > 3 {
				b.setNextLife(currentRow, currentColumn, false)
			} else if currentlyDead && liveNeighborCount == 3 {
				b.setNextLife(currentRow, currentColumn, true)
			}
		}
	}
	b.generation++
}

func (b *bufferedUniverse) countLiveNeighbors(row, column int) int {
	deltas := []int{-1, 0, 1}
	liveNeighborCount := 0
	for rowDelta := range deltas {
		for columnDelta := range deltas {
			if rowDelta == 0 && columnDelta == 0 {
				continue
			}
			if b.IsAlive(row+rowDelta, column+columnDelta) {
				liveNeighborCount++
			}
		}
	}
	return liveNeighborCount
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

func (b *bufferedUniverse) String() string {
	buffer := bytes.Buffer{}
	for currentRow := int(0); currentRow < b.RowCount(); currentRow++ {
		for currentColumns := int(0); currentColumns < b.ColumnCount(); currentColumns++ {
			if b.IsAlive(currentColumns, currentRow) {
				buffer.Write([]byte("■"))
			} else {
				buffer.Write([]byte("□"))
			}
		}
		buffer.Write([]byte("\n"))
	}
	return buffer.String()
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
