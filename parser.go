package jmatch

import "fmt"

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

func (p *parser) parseObject() {
	currentToken, nextToken := p.tokens.current(), p.tokens.next()
	if currentToken.IsRightBrace() {
		p.context = p.stack.pop()
	} else if currentToken.IsString() && nextToken.IsColon() { // key found
		p.context.setLastKey(currentToken.Value)
	} else if (currentToken.IsLeftBrace() || currentToken.IsComma()) && nextToken.IsString() {
		// pass, covered in the previous case
	} else if p.isValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBrace()) {
		// pass, will catch later values by matching them against ':'
	} else if currentToken.IsColon() {
		if p.isValue(nextToken) {
			p.matcher.Match(p.context.lastKey, nextToken)
			p.context.resetLastKey()
		} else if nextToken.IsLeftBrace() {
			p.stack.push(p.context)
			p.context = newParsingContext(p.context.lastKey, Object)
		} else if nextToken.IsLeftBracket() {
			newContextPath := p.context.lastKey
			p.context.resetLastKey()
			p.stack.push(p.context)
			p.context = newParsingContext(newContextPath, Array)
		} else {
			p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
		}
	} else {
		p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
	}
}

func (p *parser) parseArray() {
	currentToken, nextToken := p.tokens.current(), p.tokens.next()
	if currentToken.IsRightBracket() {
		p.context = p.stack.pop()
	} else if currentToken.IsLeftBracket() || currentToken.IsComma() {
		if p.isValue(nextToken) {
			p.matcher.Match(p.context.arrayPath(), nextToken)
			p.context.increaseElemsCount()
		} else if nextToken.IsLeftBracket() {
			newContextPath := p.context.arrayPath()
			p.context.increaseElemsCount()
			p.stack.push(p.context)
			p.context = newParsingContext(newContextPath, Array)
		} else if nextToken.IsLeftBrace() {
			newContextPath := p.context.arrayPath()
			p.context.increaseElemsCount()
			p.stack.push(p.context)
			p.context = newParsingContext(newContextPath, Object)
		} else {
			p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
		}
	} else if p.isValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBracket()) {
		// pass, already parsed values, arrays and object when they are after comma or left bracket
	} else {
		p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
	}
}

func (p *parser) parse() error {
	p.context = newParsingContext("", Object)

	if p.tokens.tokensCount < 2 || !p.tokens.current().IsLeftBrace() {
		p.err = fmt.Errorf("invalid JSON. must start with { and end with }")
	}

	for p.tokens.hasNext() {
		if p.context.inObject() {
			p.parseObject()
		} else if p.context.inArray() {
			p.parseArray()
		} else {
			p.err = fmt.Errorf("invalid JSON")
		}

		if p.err != nil {
			return p.err
		}

		p.tokens.move()
	}
	if !p.tokens.next().IsRightBrace() {
		p.err = fmt.Errorf("invalid JSON. the last token must be }")
	}

	return p.err
}
