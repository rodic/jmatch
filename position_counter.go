package jmatch

type positionCounter struct {
	line   int
	column int
}

func (p *positionCounter) increaseLine() {
	p.line++
	p.column = 0
}

func (p *positionCounter) increaseColumn() {
	p.column++
}

func (p *positionCounter) decreaseColumn() {
	p.column--
}

func (p *positionCounter) update(r rune) {
	if r == rune('\n') {
		p.increaseLine()
	} else {
		p.increaseColumn()
	}
}

func newRunePositionCounter() positionCounter {
	return positionCounter{
		line:   1,
		column: 0,
	}
}
