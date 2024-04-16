package parser

import t "github.com/rodic/jmatch/tokenizer"

type parenCounter struct {
	bracketCounter int
	bracesCounter  int
}

func (p *parenCounter) update(t t.Token) {
	if t.IsLeftBrace() {
		p.bracesCounter++
	} else if t.IsRightBrace() {
		p.bracesCounter--
	} else if t.IsLeftBracket() {
		p.bracketCounter++
	} else if t.IsRightBracket() {
		p.bracketCounter--
	}
}

func (p *parenCounter) isBalanced() bool {
	return p.bracesCounter == 0 && p.bracketCounter == 0
}

func newParenCounter() parenCounter {
	return parenCounter{
		bracketCounter: 0,
		bracesCounter:  0,
	}
}
