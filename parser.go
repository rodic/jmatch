package jmatch

import "fmt"

type parsingType int

const (
	Object parsingType = iota
	Array
)

type parsingContext struct {
	path        string
	elemsCount  int
	parsingType parsingType
}

func (p parsingContext) inArray() bool {
	return p.parsingType == Array
}

func (p parsingContext) inObject() bool {
	return p.parsingType == Object
}

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

type parsingResult map[string]Token

type parser struct {
	tokens         tokensList
	parsingContext parsingContext
	stack          stack
	lastKey        string
	result         parsingResult
	err            error
}

func newParser(tokens []Token) parser {
	return parser{
		tokens: newTokens(tokens),
		stack:  newStack(),
		result: make(parsingResult),
	}
}

func (p *parser) isValue(t Token) bool {
	return t.IsString() || t.IsNumber() || t.IsBoolean() || t.IsNull()
}

func (p *parser) parseObject() {
	currentToken, nextToken := p.tokens.current(), p.tokens.next()
	if currentToken.IsRightBrace() {
		p.parsingContext = p.stack.pop()
	} else if currentToken.IsString() && nextToken.IsColon() { // key found
		p.lastKey = fmt.Sprintf("%s.%s", p.lastKey, currentToken.Value)
	} else if currentToken.IsLeftBrace() && nextToken.IsString() {
		// pass, will catch later new object start with ': {' match
	} else if p.isValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBrace()) {
		// pass, will catch later values with ': v
	} else if currentToken.IsComma() && nextToken.IsString() {
		// pass, it's only valid to have string after comma
	} else if currentToken.IsColon() {
		if p.isValue(nextToken) {
			p.result[p.lastKey] = nextToken
			p.lastKey = p.parsingContext.path
		} else if nextToken.IsLeftBrace() {
			p.stack.push(p.parsingContext)
			p.parsingContext = parsingContext{
				path:        p.lastKey,
				parsingType: Object,
			}
		} else if nextToken.IsLeftBracket() {
			p.stack.push(p.parsingContext)
			p.parsingContext = parsingContext{
				path:        p.lastKey,
				parsingType: Array,
				elemsCount:  0,
			}
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
		p.parsingContext = p.stack.pop()
		p.lastKey = p.parsingContext.path
	} else if p.isValue(currentToken) && (nextToken.IsComma() || nextToken.IsRightBracket()) {
		newPath := fmt.Sprintf("%s.[%d]", p.parsingContext.path, p.parsingContext.elemsCount)
		p.parsingContext.elemsCount++
		p.result[newPath] = currentToken
	} else if nextToken.IsLeftBracket() {
		newPath := fmt.Sprintf("%s.[%d]", p.parsingContext.path, p.parsingContext.elemsCount)
		p.parsingContext.elemsCount++
		p.stack.push(p.parsingContext)
		p.parsingContext = parsingContext{
			path:        newPath,
			parsingType: Array,
			elemsCount:  0,
		}
	} else if currentToken.IsLeftBracket() {
		// skipping the first '[', all else covered in the previous case.
	} else if currentToken.IsComma() &&
		!(nextToken.IsComma() || nextToken.IsRightBrace() || nextToken.IsRightBracket()) {
		//
	} else if currentToken.IsLeftBrace() {
		p.lastKey = fmt.Sprintf("%s.[%d]", p.parsingContext.path, p.parsingContext.elemsCount)
		p.parsingContext.elemsCount++
		p.stack.push(p.parsingContext)
		p.parsingContext = parsingContext{
			path:        p.lastKey,
			parsingType: Object,
		}
	} else {
		p.err = fmt.Errorf("invalid JSON %s -> %s", currentToken.Value, nextToken.Value)
	}
}

func (p *parser) parse() (parsingResult, error) {
	p.parsingContext = parsingContext{path: "", parsingType: Object}
	p.lastKey = ""

	if p.tokens.tokensCount < 2 || !p.tokens.current().IsLeftBrace() {
		p.err = fmt.Errorf("invalid JSON. must start with { and end with }")
	}

	for p.tokens.hasNext() {
		if p.parsingContext.inObject() {
			p.parseObject()
		} else if p.parsingContext.inArray() {
			p.parseArray()
		} else {
			p.err = fmt.Errorf("invalid JSON")
		}

		if p.err != nil {
			return p.result, p.err
		}

		p.tokens.move()
	}
	if !p.tokens.next().IsRightBrace() {
		p.err = fmt.Errorf("invalid JSON. the last token must be }")
	}

	return p.result, p.err
}
