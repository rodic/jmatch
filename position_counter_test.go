package jmatch

import "testing"

func TestPositionCounter(t *testing.T) {

	c := newRunePositionCounter()

	if c.line != 1 {
		t.Errorf("Invalid initial line count")
	}

	if c.column != 0 {
		t.Error("Invalid initial column count")
	}

	c.update('a')

	if c.line != 1 {
		t.Errorf("invalid line count, expected 1, got %d", c.line)
	}

	if c.column != 1 {
		t.Errorf("invalid column count, expected 1, got %d", c.column)
	}

	c.update('ü')

	if c.line != 1 {
		t.Errorf("invalid line count, expected 1, got %d", c.line)
	}

	if c.column != 2 {
		t.Errorf("invalid column count, expected 2, got %d", c.column)
	}

	c.update('\n')

	if c.line != 2 {
		t.Errorf("invalid line count, expected 2, got %d", c.line)
	}

	if c.column != 0 {
		t.Errorf("invalid column count, expected 0, got %d", c.column)
	}

	c.update('a')

	if c.line != 2 {
		t.Errorf("invalid line count, expected 2, got %d", c.line)
	}

	if c.column != 1 {
		t.Errorf("invalid column count, expected 1, got %d", c.column)
	}

	c.decreaseColumn()

	if c.line != 2 {
		t.Errorf("invalid line count, expected 2, got %d", c.line)
	}

	if c.column != 0 {
		t.Errorf("invalid column count, expected 0, got %d", c.column)
	}
}
