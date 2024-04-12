package parser

import (
	c "github.com/rodic/jmatch/common"
	m "github.com/rodic/jmatch/matcher"
	t "github.com/rodic/jmatch/token"
)

type parser struct {
	tokens  t.TokenList
	context context
	stack   contextStack
	matcher m.Matcher
}

func NewParser(tokens []t.Token, matcher m.Matcher) parser {
	return parser{
		tokens:  t.NewTokens(tokens),
		stack:   newContextStack(),
		matcher: matcher,
	}
}

func (p *parser) isValue(t t.Token) bool {
	return t.IsString() || t.IsNumber() || t.IsBoolean() || t.IsNull()
}

func (p *parser) switchParsingContext() error {
	if p.stack.isEmpty() {
		return c.UnexpectedEndOfInputErr{}
	}

	p.context = p.stack.pop()
	return nil
}

func (p *parser) parseObject() error {
	current := p.tokens.Current()

	if current.IsRightBrace() {
		p.switchParsingContext()
		return nil
	}

	next := p.tokens.Next()

	if current.IsLeftBrace() && next.IsRightBrace() {
		return nil // pass
	}
	if current.IsComma() && p.context.isKeySet() {
		return current.AsUnexpectedTokenErr()
	}

	if current.IsColon() && !p.context.isKeySet() {
		return current.AsUnexpectedTokenErr()
	}

	if current.IsLeftBrace() || current.IsComma() {
		if next.IsString() {
			p.context.setKey(next.Value)
			p.tokens.Move()
			return nil
		} else {
			return next.AsUnexpectedTokenErr()
		}
	}

	if current.IsColon() {
		path := p.context.getPath()
		p.context.setValue()

		if p.isValue(next) {
			p.matcher.Match(path, next)
			p.tokens.Move()
		} else if next.IsLeftBrace() {
			p.stack.push(p.context)
			p.context = newObjectContext(path)
		} else if next.IsLeftBracket() {
			p.stack.push(p.context)
			p.context = newArrayContext(path)
		} else {
			return next.AsUnexpectedTokenErr()
		}

		return nil
	}
	return next.AsUnexpectedTokenErr()
}

func (p *parser) parseArray() error {
	current := p.tokens.Current()

	if current.IsRightBracket() {
		p.switchParsingContext()
		return nil
	}

	next := p.tokens.Next()

	if current.IsLeftBracket() && next.IsRightBracket() {
		return nil // pass
	}

	if current.IsLeftBracket() || current.IsComma() {
		path := p.context.getPath()
		p.context.setValue()

		if p.isValue(next) {
			p.matcher.Match(path, next)
			p.tokens.Move()
		} else if next.IsLeftBracket() {
			p.stack.push(p.context)
			p.context = newArrayContext(path)
		} else if next.IsLeftBrace() {
			p.stack.push(p.context)
			p.context = newObjectContext(path)
		} else {
			return next.AsUnexpectedTokenErr()
		}
		return nil
	}
	return current.AsUnexpectedTokenErr()
}

func (p *parser) parseContext() error {
	parenCounter := newParenCounter()

	for p.tokens.HasNext() {
		parenCounter.update(p.tokens.Current())

		if p.context.isObject() {
			err := p.parseObject()
			if err != nil {
				return err
			}
		} else if p.context.isArray() {
			err := p.parseArray()
			if err != nil {
				return err
			}
		}
		p.tokens.Move()
	}

	last := p.tokens.Current()
	parenCounter.update(last)

	if !parenCounter.isBalanced() {
		return c.UnexpectedEndOfInputErr{}
	}

	return nil
}

func (p *parser) Parse() error {

	if p.tokens.Empty() {
		return c.UnexpectedEndOfInputErr{}
	}

	first := p.tokens.Current()

	if first.IsLeftBrace() {
		p.context = newObjectContext("")
		return p.parseContext()
	}

	if first.IsLeftBracket() {
		p.context = newArrayContext("")
		return p.parseContext()
	}

	if p.isValue(first) && !p.tokens.HasNext() {
		p.matcher.Match(".", first)
		return nil
	}

	return c.UnexpectedEndOfInputErr{}
}
