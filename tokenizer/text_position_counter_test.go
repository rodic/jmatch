package tokenizer

import "testing"

func TestTextPositionCounter(t *testing.T) {

	c := newTextPositionCounter()

	if c.line != 1 {
		t.Errorf("Invalid initial line count")
	}

	if c.column != 0 {
		t.Error("Invalid initial column count")
	}

	c.increase('a')

	if c.line != 1 {
		t.Errorf("invalid line count, expected 1, got %d", c.line)
	}

	if c.column != 1 {
		t.Errorf("invalid column count, expected 1, got %d", c.column)
	}

	c.increase('ü')

	if c.line != 1 {
		t.Errorf("invalid line count, expected 1, got %d", c.line)
	}

	if c.column != 2 {
		t.Errorf("invalid column count, expected 2, got %d", c.column)
	}

	c.increase('\n')

	if c.line != 2 {
		t.Errorf("invalid line count, expected 2, got %d", c.line)
	}

	if c.column != 0 {
		t.Errorf("invalid column count, expected 0, got %d", c.column)
	}

	c.decrease('\n')

	if c.line != 1 {
		t.Errorf("invalid line count, expected 1, got %d", c.line)
	}

	if c.column != 2 {
		t.Errorf("invalid column count, expected 2, got %d", c.column)
	}

	c.increase('\n')
	c.increase('a')

	if c.line != 2 {
		t.Errorf("invalid line count, expected 2, got %d", c.line)
	}

	if c.column != 1 {
		t.Errorf("invalid column count, expected 1, got %d", c.column)
	}

	c.decrease('a')

	if c.line != 2 {
		t.Errorf("invalid line count, expected 2, got %d", c.line)
	}

	if c.column != 0 {
		t.Errorf("invalid column count, expected 0, got %d", c.column)
	}
}
