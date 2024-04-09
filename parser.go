package jmatch

import "fmt"

type stack struct {
	_stack []parsingContext
	cnt    int
}

func (s *stack) pop() parsingContext {
	s.cnt--
	top := s._stack[s.cnt]
	s._stack = s._stack[:s.cnt]
	return top
}

func (s *stack) push(stackFame parsingContext) {
	s._stack = append(s._stack, stackFame)
	s.cnt++
}

func newStack() stack {
	return stack{
		_stack: []parsingContext{},
		cnt:    0,
	}
}

type parser struct {
	tokens  tokensList
	context parsingContext
	stack   stack
	matcher Matcher
	err     error
}

func newParser(tokens []Token, matcher Matcher) parser {
	return parser{
		tokens:  newTokens(tokens),
		stack:   newStack(),
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
	} else if currentToken.IsLeftBrace() && nextToken.IsString() {
		// pass, will catch later new object start with ': {' match
	} else if p.isValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBrace()) {
		// pass, will catch later values with ': v
	} else if currentToken.IsComma() && nextToken.IsString() {
		// pass, it's only valid to have string after comma
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
	} else if p.isValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBracket()) {
		p.matcher.Match(p.context.arrayPath(), currentToken)
		p.context.increaseElemsCount()
	} else if nextToken.IsLeftBracket() {
		newContextPath := p.context.arrayPath()
		p.context.increaseElemsCount()
		p.stack.push(p.context)
		p.context = newParsingContext(newContextPath, Array)
	} else if currentToken.IsLeftBracket() {
		// skipping the first '[', all else covered in the previous case.
	} else if currentToken.IsComma() &&
		!(nextToken.IsComma() || nextToken.IsRightBrace() || nextToken.IsRightBracket()) {
		//
	} else if currentToken.IsLeftBrace() {
		newPath := p.context.arrayPath()
		p.context.increaseElemsCount()
		p.stack.push(p.context)
		p.context = newParsingContext(newPath, Object)
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
