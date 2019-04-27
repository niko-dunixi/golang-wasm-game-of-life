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
	// Single Corner:
	//   0000
	//   0000
	//   0000
	//   000X
	u := NewBufferedUniverse(4, 4, func(b *[]bool) {
		bools := *b
		bools[len(bools)-1] = true
	})
	for currentRow := 0; currentRow < u.RowCount(); currentRow++ {
		for currentColumn := 0; currentColumn < u.ColumnCount(); currentColumn++ {
			isLastCell := currentColumn+1 == u.ColumnCount() && currentRow+1 == u.RowCount()
			if isLastCell {
				if u.IsDead(currentRow, currentColumn) {
					t.Fatalf("Corner cell should be alive\n%s", u)
				}
			} else {
				if u.IsAlive(currentRow, currentColumn) {
					t.Fatalf("Only far corner cell should be alive\n%s", u)
				}
			}
		}
	}
	if u.IsDead(-1, -1) {
		t.Errorf("Didn't calculate opposite corner correctly\n%s", u)
	}
}

func TestBufferedUniverse_Iterate(t *testing.T) {
	// Still life:
	//   0000
	//   0XX0
	//   0XX0
	//   0000
	stillLife := NewBufferedUniverse(4, 4, func(b *[]bool) {
		bools := *b
		bools[5] = true
		bools[6] = true
		bools[9] = true
		bools[10] = true
	})
	if stillLife.IsDead(1, 1) || stillLife.IsDead(1, 2) ||
		stillLife.IsDead(2, 1) || stillLife.IsDead(2, 2) {
		t.Errorf("StillLife didn't initialize correctly\n%s", stillLife)
	}
	stillLife.Iterate()
	if stillLife.IsDead(1, 1) || stillLife.IsDead(1, 2) ||
		stillLife.IsDead(2, 1) || stillLife.IsDead(2, 2) {
		t.Errorf("StillLife didn't iterate correctly\n%s", stillLife)
	}
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
	oscillatorI := NewBufferedUniverse(5, 5, func(b *[]bool) {
		bools := *b
		bools[7] = true
		bools[12] = true
		bools[17] = true
	})
	if oscillatorI.IsDead(1, 2) || oscillatorI.IsDead(2, 2) || oscillatorI.IsDead(3, 2) {
		t.Errorf("Oscillator I didn't initialize correctly\n%s", oscillatorI)
	}
	oscillatorI.Iterate()
	if oscillatorI.IsDead(2, 1) || oscillatorI.IsDead(2, 2) || oscillatorI.IsDead(2, 3) {
		t.Errorf("Oscillator I didn't iterate correctly\n%s", oscillatorI)
	}
}
