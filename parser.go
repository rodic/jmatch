package jmatch

import (
	"fmt"
)

type parser struct {
	tokens  tokensList
	context context
	stack   contextStack
	matcher Matcher
	err     error
}

func newParser(tokens []Token, matcher Matcher) parser {
	return parser{
		tokens:  newTokens(tokens),
		stack:   newContextStack(),
		matcher: matcher,
	}
}

func (p *parser) isValue(t Token) bool {
	return t.IsString() || t.IsNumber() || t.IsBoolean() || t.IsNull()
}

func (p *parser) setUnexpectedEndOfInputErr() {
	p.err = fmt.Errorf("invalid JSON. Unexpected end of JSON input")
}

func (p *parser) switchParsingContext() {
	if p.stack.isEmpty() {
		p.setUnexpectedEndOfInputErr()
	} else {
		p.context = p.stack.pop()
	}
}

func (p *parser) parseObject() {
	current := p.tokens.current()

	if current.IsRightBrace() {
		p.switchParsingContext()
		return
	}

	next := p.tokens.next()

	if current.IsLeftBrace() && next.IsRightBrace() {
		// pass
	} else if current.IsComma() && !p.context.isValueSet() {
		p.err = current.toError()
	} else if current.IsColon() && p.context.isValueSet() {
		p.err = current.toError()
	} else if current.IsLeftBrace() || current.IsComma() {
		if next.IsString() && !p.context.isKeySet() {
			p.context.setKey(next.Value)
			p.tokens.move()
		} else {
			p.err = next.toError()
		}
	} else if current.IsColon() {
		path := p.context.getPath()
		p.context.setValue()

		if p.isValue(next) {
			p.matcher.Match(path, next)
			p.tokens.move()
		} else if next.IsLeftBrace() {
			p.stack.push(p.context)
			p.context = newObjectContext(path)
		} else if next.IsLeftBracket() {
			p.stack.push(p.context)
			p.context = newArrayContext(path)
		} else {
			p.err = next.toError()
		}
	} else {
		p.err = next.toError()
	}
}

func (p *parser) parseArray() {
	current := p.tokens.current()

	if current.IsRightBracket() {
		p.switchParsingContext()
		return
	}

	next := p.tokens.next()

	if current.IsLeftBracket() && next.IsRightBracket() {
		// pass
	} else if current.IsLeftBracket() || current.IsComma() {

		path := p.context.getPath()
		p.context.setValue()

		if p.isValue(next) {
			p.matcher.Match(path, next)
			p.tokens.move()
		} else if next.IsLeftBracket() {
			p.stack.push(p.context)
			p.context = newArrayContext(path)
		} else if next.IsLeftBrace() {
			p.stack.push(p.context)
			p.context = newObjectContext(path)
		} else {
			p.err = next.toError()
		}
	} else {
		p.err = current.toError()

	}
}

func (p *parser) parseSingleton() {
	current := p.tokens.current()
	if !p.isValue(current) || p.tokens.hasNext() {
		p.setUnexpectedEndOfInputErr()
	} else {
		p.matcher.Match(".", current)
	}
}

func (p *parser) parse() error {

	if p.tokens.empty() {
		p.setUnexpectedEndOfInputErr()
		return p.err
	}

	first := p.tokens.current()

	if first.IsLeftBrace() {
		p.context = newObjectContext("")
	} else if first.IsLeftBracket() {
		p.context = newArrayContext("")
	} else {
		p.parseSingleton()
		return p.err
	}

	parenCounter := newParenCounter()

	for p.tokens.hasNext() {
		parenCounter.update(p.tokens.current())

		if p.context.isObject() {
			p.parseObject()
		} else if p.context.isArray() {
			p.parseArray()
		}

		if p.err != nil {
			return p.err
		}

		p.tokens.move()
	}

	last := p.tokens.current()
	parenCounter.update(last)

	if !parenCounter.isBalanced() {
		p.setUnexpectedEndOfInputErr()
	}

	return p.err
}
