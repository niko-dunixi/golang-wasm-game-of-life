package main

import (
	"math/rand"
	"testing"
	"time"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func TestNewBufferedUniverse(t *testing.T) {
	expectedRows := r.Intn(30)
	expectedColumns := r.Intn(30)
	u := NewBufferedUniverse(expectedRows, expectedColumns, func(b *[]bool) {})
	if u.RowCount() != expectedRows {
		t.Errorf("expected rows should be %d, but was %d", expectedRows, u.RowCount())
	}
	if u.ColumnCount() != expectedColumns {
		t.Errorf("expected columns should be %d, but was %d", expectedColumns, u.ColumnCount())
	}
}

func TestBufferedUniverse_IsAlive(t *testing.T) {
	// Still life:
	//   0000
	//   0XX0
	//   0XX0
	//   0000
	u := NewBufferedUniverse(4, 4, func(b *[]bool) {
		bools := *b
		bools[5] = true
		bools[6] = true
		bools[9] = true
		bools[10] = true
	})
	if u.IsDead(1, 1) || u.IsDead(1, 2) ||
		u.IsDead(2, 1) || u.IsDead(2, 2) {
		t.Errorf("Cells weren't alive\n%s", u)
	}
}

func TestBufferedUniverse_Iterate(t *testing.T) {
	// Oscillator Ia:
	//   00000
	//   00X00
	//   00X00
	//   00X00
	//   00000
	// Oscillator Ib:
	//   00000
	//   00000
	//   0XXX0
	//   00000
	//   00000
	oscillatorOne := NewBufferedUniverse(5, 5, func(b *[]bool) {
		bools := *b
		bools[7] = true
		bools[12] = true
		bools[17] = true
	})
	if oscillatorOne.IsDead(1, 2) || oscillatorOne.IsDead(2, 2) || oscillatorOne.IsDead(3, 2) {
		t.Errorf("Cells weren't alive\n%s", oscillatorOne)
	}
}
