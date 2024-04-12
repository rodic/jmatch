package tokenizer

type textPositionCounter struct {
	line   int
	column int
}

func (p *textPositionCounter) increaseLine() {
	p.line++
	p.column = 0
}

func (p *textPositionCounter) increaseColumn() {
	p.column++
}

func (p *textPositionCounter) decreaseColumn() {
	p.column--
}

func (p *textPositionCounter) update(r rune) {
	if r == '\n' {
		p.increaseLine()
	} else {
		p.increaseColumn()
	}
}

func newTextPositionCounter() textPositionCounter {
	return textPositionCounter{
		line:   1,
		column: 0,
	}
}
