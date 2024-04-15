package parser

import (
	c "github.com/rodic/jmatch/common"
	t "github.com/rodic/jmatch/token"
	z "github.com/rodic/jmatch/tokenizer"
)

type ParsingResult struct {
	Path  string
	Token t.Token
	Error error
}

type parser struct {
	tokens       tokenList
	context      context
	stack        contextStack
	resultStream chan ParsingResult
}

func NewParser(tokenStream <-chan z.TokenResult) (*parser, error) {
	tokens, err := NewTokens(tokenStream)

	if err != nil {
		return nil, err
	}

	parser := parser{
		tokens:       *tokens,
		stack:        newContextStack(),
		resultStream: make(chan ParsingResult),
	}

	return &parser, nil
}

func (p *parser) GetResultReadStream() <-chan ParsingResult {
	return p.resultStream
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
	current := p.tokens.current

	if current.IsRightBrace() {
		p.switchParsingContext()
		return nil
	}

	next := p.tokens.next

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
			return p.tokens.move()
		} else {
			return next.AsUnexpectedTokenErr()
		}
	}
	if current.IsColon() {
		path := p.context.getPath()
		p.context.setValue()

		if p.isValue(next) {
			p.resultStream <- ParsingResult{Path: path, Token: next}
			return p.tokens.move()
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
	current := p.tokens.current

	if current.IsRightBracket() {
		p.switchParsingContext()
		return nil
	}

	next := p.tokens.next

	if current.IsLeftBracket() && next.IsRightBracket() {
		return nil // pass
	}
	if current.IsLeftBracket() || current.IsComma() {
		path := p.context.getPath()
		p.context.setValue()

		if p.isValue(next) {
			p.resultStream <- ParsingResult{Path: path, Token: next}
			return p.tokens.move()
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
	var err error

	for p.tokens.hasNext {
		parenCounter.update(p.tokens.current)

		if p.context.isObject() {
			err = p.parseObject()
		} else if p.context.isArray() {
			err = p.parseArray()
		}

		if err != nil {
			return err
		}

		err = p.tokens.move()

		if err != nil {
			return err
		}
	}

	last := p.tokens.current
	parenCounter.update(last)

	if !parenCounter.isBalanced() {
		return c.UnexpectedEndOfInputErr{}
	}

	return nil
}

func (p *parser) Parse() {
	defer close(p.resultStream)

	var err error

	first := p.tokens.current

	if p.isValue(first) && !p.tokens.hasNext {
		p.resultStream <- ParsingResult{Path: ".", Token: first}
		return
	}

	if !(first.IsLeftBrace() || first.IsLeftBracket()) {
		p.resultStream <- ParsingResult{Error: c.UnexpectedEndOfInputErr{}}
		return
	}

	if first.IsLeftBrace() {
		p.context = newObjectContext("")
	}

	if first.IsLeftBracket() {
		p.context = newArrayContext(".")
	}

	err = p.parseContext()

	if err != nil {
		p.resultStream <- ParsingResult{Error: err}
		return
	}

}
