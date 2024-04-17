package tokenizer

type textPositionCounter struct {
	line           int
	column         int
	lastLineColumn int
}

func (p *textPositionCounter) increase(r rune) {
	if r == '\n' {
		p.line++
		p.lastLineColumn = p.column
		p.column = 0
	} else {
		p.column++
	}
}

func (p *textPositionCounter) decrease(r rune) {
	if r == '\n' {
		p.line--
		p.column = p.lastLineColumn
	} else {
		p.column--
	}
}

func newTextPositionCounter() textPositionCounter {
	return textPositionCounter{
		line:           1,
		column:         0,
		lastLineColumn: 0,
	}
}
