package jmatch

import (
	"fmt"
)

type parser struct {
	tokens  tokensList
	context parsingContext
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

func (p *parser) parseObject() {
	if p.tokens.current().IsRightBrace() {
		p.context = p.stack.pop() // switch to previous context
	} else {
		current, next := p.tokens.current(), p.tokens.next()

		if current.IsLeftBrace() && next.IsRightBrace() {
			// pass
		} else if current.IsComma() && !p.context.isValueSet() {
			p.err = current.toError()
		} else if current.IsColon() && p.context.isValueSet() {
			p.err = current.toError()
		} else if current.IsLeftBrace() || current.IsComma() {
			if next.IsString() && !p.context.isKeySet() {
				p.context.setKey(p.tokens.next().Value)
				p.tokens.move()
			} else {
				p.err = next.toError()
			}
		} else if current.IsColon() && p.context.isKeySet() {

			path := p.context.getPath()
			p.context.setValue()

			if p.isValue(next) {
				p.matcher.Match(path, next)
				p.tokens.move()
			} else if next.IsLeftBrace() {
				p.stack.push(p.context)
				p.context = newParsingContext(path, Object)
			} else if next.IsLeftBracket() {
				p.stack.push(p.context)
				p.context = newParsingContext(path, Array)
			} else {
				p.err = next.toError()
			}

		} else {
			p.err = next.toError()
		}
	}
}

func (p *parser) parseArray() {
	if p.tokens.current().IsRightBracket() {
		p.context = p.stack.pop() // switch to previous context
	} else {
		current, next := p.tokens.current(), p.tokens.next()

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
				p.context = newParsingContext(path, Array)
			} else if next.IsLeftBrace() {
				p.stack.push(p.context)
				p.context = newParsingContext(path, Object)
			} else {
				p.err = next.toError()

			}
		} else {
			p.err = current.toError()

		}
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
		p.context = newParsingContext("", Object)
	} else if first.IsLeftBracket() {
		p.context = newParsingContext("", Array)
	} else {
		p.parseSingleton()
		return p.err
	}

	parenCounter := newParenCounter()

	for p.tokens.hasNext() {
		parenCounter.update(p.tokens.current())

		if p.context.inObject() {
			p.parseObject()
		} else if p.context.inArray() {
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
