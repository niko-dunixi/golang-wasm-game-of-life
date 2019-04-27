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
	u := NewBufferedUniverse(expectedRows, expectedColumns, NewRand(1))
	if u.RowCount() != expectedRows {
		t.Errorf("expected rows should be %d, but was %d", expectedRows, u.RowCount())
	}
	if u.ColumnCount() != expectedColumns {
		t.Errorf("expected columns should be %d, but was %d", expectedColumns, u.ColumnCount())
	}
}
